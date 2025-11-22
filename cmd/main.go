package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/konnen/review-assign-service/internal/config"
	"github.com/konnen/review-assign-service/internal/database/postgres"
	"github.com/konnen/review-assign-service/internal/handler"
	"github.com/konnen/review-assign-service/internal/middleware/logger"
	"github.com/konnen/review-assign-service/internal/repository"
	"github.com/konnen/review-assign-service/internal/service"
	"github.com/konnen/review-assign-service/internal/validator"
)

func main() {
	lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	cfg, err := config.LoadConfig()
	if err != nil {
		lgr.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	lgr.Info("config loaded", "cfg", cfg)

	db, err := postgres.OpenDB(cfg.DBUrl)
	defer func(){
		if err := db.Close(); err != nil{
        	lgr.Error("failed to close database", "error", err)
		}
	}()
	if err != nil {
		lgr.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	v := validator.NewValidator()
	r := repository.NewRepository(db)
	s := service.NewService(trManager, r.Team, r.User, r.PullRequest)
	h := handler.NewHandler(lgr, v, s.Team, s.User, s.PullRequest)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(logger.New(lgr))

	h.InitRoutes(router)

	srv := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}

	lgr.Info("starting server ...")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lgr.Error("listen", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	lgr.Info("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		lgr.Error("server shutdown:", "error", err)
		os.Exit(1)
	}
	<-ctx.Done()
	lgr.Info("timeout of 5 seconds.")
	lgr.Info("server exiting")
}

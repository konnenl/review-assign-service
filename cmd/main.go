package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/konnen/review-assign-service/internal/config"
	"github.com/konnen/review-assign-service/internal/handler"
	"github.com/konnen/review-assign-service/internal/middleware/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	cfg, err := config.LoadConfig()
	if err != nil {
		lgr.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	lgr.Info("config loaded", "cfg", cfg)
	h := handler.New(lgr)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(logger.New(lgr))

	h.InitRoutes(r)

	srv := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: r,
	}

	lgr.Info("starting Server ...")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lgr.Error("listen: %s\\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	lgr.Info("shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		lgr.Error("server shutdown:", err)
		os.Exit(1)
	}
	select {
	case <-ctx.Done():
		lgr.Info("timeout of 5 seconds.")
	}
	lgr.Info("server exiting")
}

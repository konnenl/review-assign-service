package main

import (
	"log/slog"
	"github.com/konnen/review-assign-service/internal/config"
	"os"
)

func main() {
	lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	cfg, err := config.LoadConfig()
	if err != nil {
		lgr.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	lgr.Info("config loaded", "cfg", cfg)
}

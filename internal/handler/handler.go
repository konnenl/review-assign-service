package handler

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type Handler struct {
	lgr *slog.Logger
}

func New(lgr *slog.Logger) *Handler {
	return &Handler{
		lgr: lgr,
	}
}

func (h *Handler) InitRoutes(r *chi.Mux) *chi.Mux {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	return r
}

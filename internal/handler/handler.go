package handler

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"encoding/json"
	"github.com/konnen/review-assign-service/internal/validator"
)

type Handler struct {
	lgr  *slog.Logger
	team *teamHandler
}

func NewHandler(lgr *slog.Logger, validator *validator.CustomValidator, teamService teamService) *Handler {
	return &Handler{
		lgr:  lgr,
		team: newTeamHandler(lgr, validator, teamService),
	}
}

func (h *Handler) InitRoutes(r *chi.Mux) *chi.Mux {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.team.addTeam)
		r.Get("/get", h.team.getTeam)
	})
	return r
}

func respondWithError(w http.ResponseWriter, statusCode int, pth string, err error, message string, lgr *slog.Logger) {
	lgr.Error(err.Error(), "handler", pth)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := map[string]string{
		"error": message,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
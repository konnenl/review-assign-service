package handler

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type Handler struct {
	lgr  *slog.Logger
	team *teamHandler
}

func NewHandler(lgr *slog.Logger, teamService teamService) *Handler {
	return &Handler{
		lgr:  lgr,
		team: newTeamHandler(lgr, teamService),
	}
}

func (h *Handler) InitRoutes(r *chi.Mux) *chi.Mux {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.team.addTeam)
		r.Get("/get", h.team.getTeam)
	})
	return r
}

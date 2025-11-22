package handler

import (
	"encoding/json"
	"log/slog"
	"net/http" 

	"github.com/go-chi/chi/v5"

	"github.com/konnen/review-assign-service/internal/dto"
	"github.com/konnen/review-assign-service/internal/validator"
)

type Handler struct {
	lgr  *slog.Logger
	team *teamHandler
	user *userHandler
}

func NewHandler(lgr *slog.Logger, validator *validator.CustomValidator, teamService teamService, userService userService) *Handler {
	return &Handler{
		lgr:  lgr,
		team: newTeamHandler(lgr, validator, teamService),
		user: newUserHandler(lgr, validator, userService),
	}
}

func (h *Handler) InitRoutes(r *chi.Mux) *chi.Mux {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.team.addTeam)
		r.Get("/get", h.team.getTeam)
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", h.user.setIsActive)
	})
	return r
}

func respondWithError(w http.ResponseWriter, statusCode int, pth string, err error, resp dto.ErrorResp, lgr *slog.Logger) {
	lgr.Error(err.Error(), "handler", pth)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(resp)
}

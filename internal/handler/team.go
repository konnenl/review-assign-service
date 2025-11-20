package handler

import (
	"log/slog"
	"net/http"
)

type teamService interface {
}

type teamHandler struct {
	lgr         *slog.Logger
	teamService teamService
}

func newTeamHandler(lgr *slog.Logger, teamService teamService) *teamHandler {
	return &teamHandler{
		lgr:         lgr,
		teamService: teamService,
	}
}

func (h *teamHandler) addTeam(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("add team"))
}

func (h *teamHandler) getTeam(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get team"))
}

package handler

import (
	"log/slog"
	"net/http"
	"encoding/json"
	"github.com/konnen/review-assign-service/internal/validator"
	"github.com/konnen/review-assign-service/internal/dto"
)

type teamService interface {
}

type teamHandler struct {
	lgr         *slog.Logger
	validator *validator.CustomValidator
	teamService teamService
}

func newTeamHandler(lgr *slog.Logger, validator *validator.CustomValidator, teamService teamService) *teamHandler {
	return &teamHandler{
		lgr:         lgr,
		validator: validator,
		teamService: teamService,
	}
}

//POST /team/add
func (h *teamHandler) addTeam(w http.ResponseWriter, r *http.Request) {
	const pth = "handler.team.addTeam"
	var team dto.TeamReq
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		respondWithError(w, http.StatusBadRequest, pth, err, "invalid request body", h.lgr)
        return
	}

	if err := h.validator.Validate(&team); err != nil{
		respondWithError(w, http.StatusBadRequest, pth, validator.GetValidationErrors(err), "invalid request body", h.lgr)
        return
	}

	w.Write([]byte("add team"))
}

func (h *teamHandler) getTeam(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get team"))
}

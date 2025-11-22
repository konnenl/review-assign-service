package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/konnen/review-assign-service/internal/dto"
	"github.com/konnen/review-assign-service/internal/errs"
	"github.com/konnen/review-assign-service/internal/validator"
	"log/slog"
	"net/http"
)

type teamService interface {
	AddTeamWithMembers(ctx context.Context, team dto.Team) error
}

type teamHandler struct {
	lgr         *slog.Logger
	validator   *validator.CustomValidator
	teamService teamService
}

func newTeamHandler(lgr *slog.Logger, validator *validator.CustomValidator, teamService teamService) *teamHandler {
	return &teamHandler{
		lgr:         lgr,
		validator:   validator,
		teamService: teamService,
	}
}

// POST /team/add
func (h *teamHandler) addTeam(w http.ResponseWriter, r *http.Request) {
	const pth = "handler.team.addTeam"
	var team dto.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgValidationError
		respondWithError(w, http.StatusBadRequest, pth, err, resp, h.lgr)
		return
	}

	if err := h.validator.Validate(&team); err != nil {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgValidationError
		respondWithError(w, http.StatusBadRequest, pth, validator.GetValidationErrors(err), resp, h.lgr)
		return
	}

	if err := h.teamService.AddTeamWithMembers(r.Context(), team); err != nil {
		if errors.Is(err, errs.ErrTeamExists) {
			resp := dto.ErrorResp{}
			resp.Error.Code = dto.CodeTeamExists
			resp.Error.Message = dto.MsgValidationError
			respondWithError(w, http.StatusBadRequest, pth, err, resp, h.lgr)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := dto.TeamResp{Team: team}
	_ = json.NewEncoder(w).Encode(resp)

}

func (h *teamHandler) getTeam(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get team"))
}

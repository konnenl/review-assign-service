package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	
	"github.com/konnen/review-assign-service/internal/dto"
	"github.com/konnen/review-assign-service/internal/errs"
	"github.com/konnen/review-assign-service/internal/model"
	"github.com/konnen/review-assign-service/internal/validator"
)

type teamService interface {
	AddTeamWithMembers(ctx context.Context, team model.Team) error
	GetTeamWithMembers(ctx context.Context, teamName string) (model.Team, error)
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
	var team model.Team
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
		resp := dto.ErrorResp{}
		if errors.Is(err, errs.ErrTeamExists) {
			resp.Error.Code = dto.CodeTeamExists
			resp.Error.Message = dto.MsgValidationError
		} else {
			resp.Error.Code = dto.CodeInternalError
			resp.Error.Message = dto.MsgInternalError
		}
		respondWithError(w, http.StatusBadRequest, pth, err, resp, h.lgr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := dto.TeamResp{Team: team}
	_ = json.NewEncoder(w).Encode(resp)

}

// GET /team/get?team_name=...
func (h *teamHandler) getTeam(w http.ResponseWriter, r *http.Request) {
	const pth = "handler.team.getTeam"
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgInvalidQueryError
		respondWithError(w, http.StatusBadRequest, pth, fmt.Errorf("%s: missing team_name", pth), resp, h.lgr)
	}
	team, err := h.teamService.GetTeamWithMembers(r.Context(), teamName)
	if err != nil {
		resp := dto.ErrorResp{}
		if errors.Is(err, errs.ErrTeamNotFound) {
			resp.Error.Code = dto.CodeNotFound
			resp.Error.Message = dto.MsgNotFound
		} else {
			resp.Error.Code = dto.CodeInternalError
			resp.Error.Message = dto.MsgInternalError
		}
		respondWithError(w, http.StatusBadRequest, pth, err, resp, h.lgr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(team)
}

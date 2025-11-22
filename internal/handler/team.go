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
	"github.com/konnen/review-assign-service/internal/mapper"
	"github.com/konnen/review-assign-service/internal/model"
	"github.com/konnen/review-assign-service/internal/validator"
)

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
	var teamDTO dto.TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&teamDTO); err != nil {
		respondWithValidationError(w, pth, err, h.lgr)
		return
	}

	if err := h.validator.Validate(&teamDTO); err != nil {
		respondWithValidationError(w, pth, validator.GetValidationErrors(err), h.lgr)
		return
	}

	teamModel := mapper.TeamtoModel(teamDTO)
	if err := h.teamService.AddTeamWithMembers(r.Context(), teamModel); err != nil {
		resp := dto.ErrorResp{}
		if errors.Is(err, errs.ErrTeamExists) {
			resp.Error.Code = dto.CodeTeamExists
			resp.Error.Message = dto.MsgValidationError
			respondWithError(w, http.StatusBadRequest, pth, err, resp, h.lgr)
		} else {
			resp.Error.Code = dto.CodeInternalError
			resp.Error.Message = dto.MsgInternalError
			respondWithError(w, http.StatusInternalServerError, pth, err, resp, h.lgr)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := dto.TeamResp{Team: teamDTO}
	_ = json.NewEncoder(w).Encode(resp)
}

// GET /team/get?team_name=...
func (h *teamHandler) getTeam(w http.ResponseWriter, r *http.Request) {
	const pth = "handler.team.getTeam"
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		respondWithValidationError(w, pth, fmt.Errorf("%s: missing team_name", pth), h.lgr)
		return
	}
	teamModel, err := h.teamService.GetTeamWithMembers(r.Context(), teamName)
	if err != nil {
		resp := dto.ErrorResp{}
		if errors.Is(err, errs.ErrTeamNotFound) {
			resp.Error.Code = dto.CodeNotFound
			resp.Error.Message = dto.MsgNotFound
			respondWithError(w, http.StatusNotFound, pth, err, resp, h.lgr)
		} else {
			resp.Error.Code = dto.CodeInternalError
			resp.Error.Message = dto.MsgInternalError
			respondWithError(w, http.StatusInternalServerError, pth, err, resp, h.lgr)
		}
		return
	}
	teamDTO := mapper.TeamToDTO(teamModel)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(teamDTO)
}

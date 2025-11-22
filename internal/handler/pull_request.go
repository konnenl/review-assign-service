package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/konnen/review-assign-service/internal/dto"
	"github.com/konnen/review-assign-service/internal/errs"
	"github.com/konnen/review-assign-service/internal/mapper"
	"github.com/konnen/review-assign-service/internal/model"
	"github.com/konnen/review-assign-service/internal/validator"
)

type pullRequestService interface {
	CreatePullRequest(ctx context.Context, pr model.PullRequest) (model.PullRequest, error)
	Merge(ctx context.Context, prID string) (model.PullRequest, error)
	Reassign(ctx context.Context, pullRequeprID, oldReviewerID string) (model.PullRequest, string, error)
}

type pullRequestHandler struct {
	lgr                *slog.Logger
	validator          *validator.CustomValidator
	pullRequestService pullRequestService
}

func newPullRequestHandler(lgr *slog.Logger, validator *validator.CustomValidator, pullRequestService pullRequestService) *pullRequestHandler {
	return &pullRequestHandler{
		lgr:                lgr,
		validator:          validator,
		pullRequestService: pullRequestService,
	}
}

func (h *pullRequestHandler) createPullRequest(w http.ResponseWriter, r *http.Request) {
	const pth = "handler.pullRequest.createPullRequest"
	var prDTO dto.PullRequestShort
	if err := json.NewDecoder(r.Body).Decode(&prDTO); err != nil {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgValidationError
		respondWithError(w, http.StatusBadRequest, pth, err, resp, h.lgr)
		return
	}

	if err := h.validator.Validate(&prDTO); err != nil {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgValidationError
		respondWithError(w, http.StatusBadRequest, pth, validator.GetValidationErrors(err), resp, h.lgr)
		return
	}

	prShortModel := mapper.PRtoModel(prDTO)
	prShortModel.Status = model.StatusOpen
	prModel, err := h.pullRequestService.CreatePullRequest(r.Context(), prShortModel)
	if err != nil {
		resp := dto.ErrorResp{}
		if errors.Is(err, errs.ErrTeamNotFound) {
			resp.Error.Code = dto.CodeNotFound
			resp.Error.Message = dto.MsgNotFound
			respondWithError(w, http.StatusNotFound, pth, err, resp, h.lgr)
		} else if errors.Is(err, errs.ErrPRExists) {
			resp.Error.Code = dto.CodePRExists
			resp.Error.Message = dto.MsgPRExists
			respondWithError(w, http.StatusConflict, pth, err, resp, h.lgr)
		} else {
			resp.Error.Code = dto.CodeInternalError
			resp.Error.Message = dto.MsgInternalError
			respondWithError(w, http.StatusInternalServerError, pth, err, resp, h.lgr)
		}
		return
	}
	prResp := mapper.PRtoDTO(prModel)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := dto.PullRequestResp{PullRequest: prResp}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *pullRequestHandler) merge(w http.ResponseWriter, r *http.Request) {
	const pth = "handler.pullRequest.merge"
	var mergeDTO dto.Merge
	if err := json.NewDecoder(r.Body).Decode(&mergeDTO); err != nil {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgValidationError
		respondWithError(w, http.StatusBadRequest, pth, err, resp, h.lgr)
		return
	}

	if err := h.validator.Validate(&mergeDTO); err != nil {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgValidationError
		respondWithError(w, http.StatusBadRequest, pth, validator.GetValidationErrors(err), resp, h.lgr)
		return
	}

	prModel, err := h.pullRequestService.Merge(r.Context(), mergeDTO.PullRequestID)
	if err != nil {
		resp := dto.ErrorResp{}
		if errors.Is(err, errs.ErrPRNotFound) {
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
	prResp := mapper.PRtoDTO(prModel)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := dto.PullRequestResp{PullRequest: prResp}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *pullRequestHandler) reassign(w http.ResponseWriter, r *http.Request) {
	const pth = "handler.pullRequest.reassign"
	var reassignDTO dto.Reassign
	if err := json.NewDecoder(r.Body).Decode(&reassignDTO); err != nil {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgValidationError
		respondWithError(w, http.StatusBadRequest, pth, err, resp, h.lgr)
		return
	}

	if err := h.validator.Validate(&reassignDTO); err != nil {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgValidationError
		respondWithError(w, http.StatusBadRequest, pth, validator.GetValidationErrors(err), resp, h.lgr)
		return
	}

	prModel, replacedBy, err := h.pullRequestService.Reassign(r.Context(), reassignDTO.PullRequestID, reassignDTO.OldReviewerID)
	if err != nil {
		resp := dto.ErrorResp{}
		if errors.Is(err, errs.ErrPRNotFound) || errors.Is(err, errs.ErrUserNotFound) {
			resp.Error.Code = dto.CodeNotFound
			resp.Error.Message = dto.MsgNotFound
			respondWithError(w, http.StatusNotFound, pth, err, resp, h.lgr)
		} else if errors.Is(err, errs.ErrPRMerged) {
			resp.Error.Code = dto.CodePRMerged
			resp.Error.Message = dto.MsgPRMerged
			respondWithError(w, http.StatusConflict, pth, err, resp, h.lgr)
		} else {
			resp.Error.Code = dto.CodeInternalError
			resp.Error.Message = dto.MsgInternalError
			respondWithError(w, http.StatusInternalServerError, pth, err, resp, h.lgr)
		}
		return
	}

	prResp := mapper.PRtoDTO(prModel)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if replacedBy == "" {
		replacedBy = " "
	}
	resp := dto.PullRequestResp{PullRequest: prResp, ReplacedBy: replacedBy}
	_ = json.NewEncoder(w).Encode(resp)
}

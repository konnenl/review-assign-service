package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"fmt"
	
	"github.com/konnen/review-assign-service/internal/dto"
	"github.com/konnen/review-assign-service/internal/errs"
	"github.com/konnen/review-assign-service/internal/validator"
	"github.com/konnen/review-assign-service/internal/mapper"
	"github.com/konnen/review-assign-service/internal/model"
)

type userService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (model.User, error)
	GetReview(ctx context.Context, userID string) ([]model.PullRequest, error)
}

type userHandler struct {
	lgr         *slog.Logger
	validator   *validator.CustomValidator
	userService userService
}

func newUserHandler(lgr *slog.Logger, validator *validator.CustomValidator, userService userService) *userHandler {
	return &userHandler{
		lgr:         lgr,
		validator:   validator,
		userService: userService,
	}
}

func (h *userHandler) setIsActive(w http.ResponseWriter, r *http.Request) {
	const pth = "handler.user.setIsActive"
	var isActiveReq dto.SetIsActiveReq
	if err := json.NewDecoder(r.Body).Decode(&isActiveReq); err != nil {
		//TODO respondWithValidationError
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgValidationError
		respondWithError(w, http.StatusBadRequest, pth, err, resp, h.lgr)
		return
	}

	if err := h.validator.Validate(&isActiveReq); err != nil {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgValidationError
		respondWithError(w, http.StatusBadRequest, pth, validator.GetValidationErrors(err), resp, h.lgr)
		return
	}
	userModel, err := h.userService.SetIsActive(r.Context(), isActiveReq.UserID, isActiveReq.IsActive)
	if err != nil {
		resp := dto.ErrorResp{}
		if errors.Is(err, errs.ErrUserNotFound) {
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
	userWithTeamDTO := mapper.UserToWithTeamDTO(userModel)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := dto.UserResp{User: userWithTeamDTO}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *userHandler) getReview(w http.ResponseWriter, r *http.Request) {
	const pth = "handler.user.getReview"
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeValidationError
		resp.Error.Message = dto.MsgInvalidQueryError
		respondWithError(w, http.StatusBadRequest, pth, fmt.Errorf("%s: missing user_id", pth), resp, h.lgr)
		return
	}

	prsModel, err := h.userService.GetReview(r.Context(), userID)
	if err != nil {
		resp := dto.ErrorResp{}
		resp.Error.Code = dto.CodeInternalError
		resp.Error.Message = dto.MsgInternalError
		respondWithError(w, http.StatusInternalServerError, pth, err, resp, h.lgr)
		return
	}
	userReview := mapper.PRtoUserReview(userID, prsModel)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(userReview)
}
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/konnen/review-assign-service/internal/dto"
	"github.com/konnen/review-assign-service/internal/validator"
	"github.com/konnen/review-assign-service/internal/model"
)

type Handler struct {
	lgr         *slog.Logger
	team        *teamHandler
	user        *userHandler
	pullRequest *pullRequestHandler
}

type teamService interface {
	AddTeamWithMembers(ctx context.Context, team model.Team) error
	GetTeamWithMembers(ctx context.Context, teamName string) (model.Team, error)
}

type userService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (model.User, error)
	GetReview(ctx context.Context, userID string) ([]model.PullRequest, error)
}

type pullRequestService interface {
	CreatePullRequest(ctx context.Context, pr model.PullRequest) (model.PullRequest, error)
	Merge(ctx context.Context, prID string) (model.PullRequest, error)
	Reassign(ctx context.Context, pullRequeprID, oldReviewerID string) (model.PullRequest, string, error)
}

func NewHandler(lgr *slog.Logger, validator *validator.CustomValidator, teamService teamService, userService userService, pullRequestService pullRequestService) *Handler {
	return &Handler{
		lgr:         lgr,
		team:        newTeamHandler(lgr, validator, teamService),
		user:        newUserHandler(lgr, validator, userService),
		pullRequest: newPullRequestHandler(lgr, validator, pullRequestService),
	}
}

func (h *Handler) InitRoutes(r *chi.Mux) *chi.Mux {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.team.addTeam)
		r.Get("/get", h.team.getTeam)
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", h.user.setIsActive)
		r.Get("/getReview", h.user.getReview)
	})
	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", h.pullRequest.createPullRequest)
		r.Post("/merge", h.pullRequest.merge)
		r.Post("/reassign", h.pullRequest.reassign)
	})
	return r
}

func respondWithError(w http.ResponseWriter, statusCode int, pth string, err error, resp dto.ErrorResp, lgr *slog.Logger) {
	lgr.Error(err.Error(), "handler", pth)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(resp)
}

func respondWithValidationError(w http.ResponseWriter, pth string, err error, lgr *slog.Logger) {
	resp := dto.ErrorResp{}
	resp.Error.Code = dto.CodeValidationError
	resp.Error.Message = dto.MsgValidationError
	lgr.Error(err.Error(), "handler", pth)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	_ = json.NewEncoder(w).Encode(resp)
}

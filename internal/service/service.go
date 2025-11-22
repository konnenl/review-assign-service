package service

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/konnen/review-assign-service/internal/model"
)

type Service struct {
	Team        *teamService
	User        *userService
	PullRequest *pullRequestService
}

func NewService(trManager *manager.Manager, teamRepo teamRepository, userRepo userRepository, pullRequestRepo pullRequestRepository) *Service {
	return &Service{
		Team:        newTeamService(trManager, teamRepo, userRepo),
		User:        newUserService(trManager, teamRepo, userRepo, pullRequestRepo),
		PullRequest: newPullRequestService(trManager, pullRequestRepo, teamRepo),
	}
}

type teamRepository interface {
	AddTeam(ctx context.Context, team model.Team) error
	AssignUserToTeam(ctx context.Context, user model.User, teamName string) error
	IsTeamExists(ctx context.Context, name string) (bool, error)
	GetTeamByAuthor(ctx context.Context, authorID string) (model.Team, error)
	GetTeamWithMembers(ctx context.Context, teamName string) (model.Team, error)
}

type userRepository interface {
	AddUser(ctx context.Context, user model.User) error
	IsExistUser(ctx context.Context, id string) (bool, error)
	SetISActive(ctx context.Context, userID string, isActive bool) (model.User, error)
}

type pullRequestRepository interface {
	GetReviews(ctx context.Context, userID string) ([]model.PullRequest, error)
	GetByID(ctx context.Context, id string) (model.PullRequest, error)
	CreatePullRequest(ctx context.Context, pr model.PullRequest) error
	AssignReviewer(ctx context.Context, prID, userID string) error
}

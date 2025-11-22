package service

import (
	"context"
	"fmt"
	"errors"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/konnen/review-assign-service/internal/errs"
	"github.com/konnen/review-assign-service/internal/model"
)

type userService struct {
	teamRepo  teamRepository
	userRepo  userRepository
	pullRequestRepo  pullRequestRepository
	trManager *manager.Manager
}

func newUserService(trManager *manager.Manager, teamRepo teamRepository, userRepo userRepository, pullRequestRepo  pullRequestRepository) *userService {
	return &userService{
		teamRepo:  teamRepo,
		userRepo:  userRepo,
		pullRequestRepo: pullRequestRepo,
		trManager: trManager,
	}
}

func (s *userService) SetIsActive(ctx context.Context, userID string, isActive bool) (model.User, error){
	const pth = "service.user.SetIsActive"
	user, err := s.userRepo.SetISActive(ctx, userID, isActive)
	if err != nil{
		if errors.Is(err, errs.ErrUserNotFound){
			return model.User{}, err
		}
		return model.User{}, fmt.Errorf("%s: setting user is_active: %w", pth, err)
	}
	return user, nil
}

func (s *userService) GetReview(ctx context.Context, userID string) ([]model.PullRequest, error){
	const pth = "service.user.GetReview"
	prs, err := s.pullRequestRepo.GetReviews(ctx, userID)
	if err != nil{
		return nil, fmt.Errorf("%s: getting users reviews: %w", pth, err)
	}
	return prs, nil
}
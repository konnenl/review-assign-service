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
	trManager *manager.Manager
}

func newUserService(trManager *manager.Manager, teamRepo teamRepository, userRepo userRepository) *userService {
	return &userService{
		teamRepo:  teamRepo,
		userRepo:  userRepo,
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
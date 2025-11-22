package service

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/konnen/review-assign-service/internal/errs"
	"github.com/konnen/review-assign-service/internal/model"
)

type teamService struct {
	teamRepo  teamRepository
	userRepo  userRepository
	trManager *manager.Manager
}

func newTeamService(trManager *manager.Manager, teamRepo teamRepository, userRepo userRepository) *teamService {
	return &teamService{
		teamRepo:  teamRepo,
		userRepo:  userRepo,
		trManager: trManager,
	}
}

func (s *teamService) AddTeamWithMembers(ctx context.Context, team model.Team) error {
	const pth = "service.team.AddTeamWithMembers"
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		exists, err := s.teamRepo.IsTeamExists(ctx, team.Name)
		if err != nil {
			return fmt.Errorf("%s: checking team existence: %w", pth, err)
		}
		if exists {
			return errs.ErrTeamExists
		}

		if err := s.teamRepo.AddTeam(ctx, team); err != nil {
			return fmt.Errorf("%s: adding team: %w", pth, err)
		}

		for _, member := range team.Members {
			exists, err := s.userRepo.IsExistUser(ctx, member.ID)
			if err != nil {
				return fmt.Errorf("%s: checking user existence: %w", pth, err)
			}
			if !exists {
				if err := s.userRepo.AddUser(ctx, member); err != nil {
					return fmt.Errorf("%s: adding user: %w", pth, err)
				}
			}

			if err := s.teamRepo.AssignUserToTeam(ctx, member, team.Name); err != nil {
				return fmt.Errorf("%s: adding member to team: %w", pth, err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", pth, err)
	}
	return nil
}

func (s *teamService) GetTeamWithMembers(ctx context.Context, teamName string) (model.Team, error) {
	const pth = "service.team.GetTeamWithMembers"
	var team model.Team
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		exists, err := s.teamRepo.IsTeamExists(ctx, teamName)
		if err != nil {
			return fmt.Errorf("checking team existence: %w", err)
		}
		if !exists {
			return errs.ErrTeamNotFound
		}
		team, err = s.teamRepo.GetTeamWithMembers(ctx, teamName)
		if err != nil {
			return fmt.Errorf("getting team with members: %w", err)
		}
		return nil
	})
	if err != nil {
		return team, fmt.Errorf("%s: %w", pth, err)
	}
	return team, nil
}

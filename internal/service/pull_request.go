package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/konnen/review-assign-service/internal/errs"
	"github.com/konnen/review-assign-service/internal/model"
)

type pullRequestService struct {
	pullRequestRepo pullRequestRepository
	teamRepo        teamRepository
	trManager       *manager.Manager
}

func newPullRequestService(trManager *manager.Manager, pullRequestRepo pullRequestRepository, teamRepo teamRepository) *pullRequestService {
	rand.Seed(time.Now().UTC().UnixNano())
	return &pullRequestService{
		pullRequestRepo: pullRequestRepo,
		teamRepo:        teamRepo,
		trManager:       trManager,
	}
}

func (s *pullRequestService) CreatePullRequest(ctx context.Context, pr model.PullRequest) (model.PullRequest, error) {
	const pth = "service.pullRequest.CreatePullRequest"
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		team, err := s.teamRepo.GetTeamByAuthor(ctx, pr.AuthorID)
		if err != nil {
			if errors.Is(err, errs.ErrTeamNotFound) {
				return errs.ErrTeamNotFound
			}
			return fmt.Errorf("checking team existence: %w", err)
		}

		_, err = s.pullRequestRepo.GetByID(ctx, pr.ID)
		if err != nil && !errors.Is(err, errs.ErrPRNotFound) {
			return fmt.Errorf("checking existing PR: %w", err)
		}
		if err == nil {
			return errs.ErrPRExists
		}

		var members []model.User
		for _, u := range team.Members {
			if u.ID != pr.AuthorID {
				members = append(members, u)
			}
		}
		reviewers := randomUsers(members, 2)

		if err := s.pullRequestRepo.CreatePullRequest(ctx, pr); err != nil {
			return fmt.Errorf("creating PR: %w", err)
		}

		for _, u := range reviewers {
			if err := s.pullRequestRepo.AssignReviewer(ctx, pr.ID, u.ID); err != nil {
				return fmt.Errorf("assigning reviewer %s: %w", u.ID, err)
			}
		}

		pr.AssignedReviewers = reviewers
		return nil
	})
	if err != nil {
		return model.PullRequest{}, fmt.Errorf("%s: %w", pth, err)
	}
	return pr, nil
}

func randomUsers(members []model.User, n int) []model.User {
	if len(members) == 0 || n <= 0 {
		return []model.User{}
	}

	if n > len(members) {
		n = len(members)
	}

	perm := rand.Perm(len(members))
	result := make([]model.User, n)
	for i := 0; i < n; i++ {
		result[i] = members[perm[i]]
	}

	return result
}

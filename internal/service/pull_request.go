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
	userRepo        userRepository
	trManager       *manager.Manager
}

func newPullRequestService(trManager *manager.Manager, pullRequestRepo pullRequestRepository, teamRepo teamRepository, userRepo userRepository) *pullRequestService {
	rand.Seed(time.Now().UTC().UnixNano())
	return &pullRequestService{
		pullRequestRepo: pullRequestRepo,
		teamRepo:        teamRepo,
		userRepo:        userRepo,
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
			return fmt.Errorf("getting PR: %w", err)
		}
		if err == nil {
			return errs.ErrPRExists
		}

		var members []model.User
		for _, u := range team.Members {
			if u.ID != pr.AuthorID && u.IsActive != nil && *u.IsActive == true {
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

func (s *pullRequestService) Merge(ctx context.Context, prID string) (model.PullRequest, error){
	const pth = "service.pullRequest.Merge"
	var pr model.PullRequest
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		var err error
		pr, err = s.pullRequestRepo.GetByID(ctx, prID)
		if err != nil{
			if errors.Is(err, errs.ErrPRNotFound){
				return errs.ErrPRNotFound
			}
			return fmt.Errorf("getting PR: %w", err)
		}

		if pr.Status == model.StatusMerged {
			return nil
		}

		timeNow := time.Now().UTC()

		if err := s.pullRequestRepo.Merge(ctx, prID, timeNow); err != nil {
			return fmt.Errorf("merging PR: %w", err)
		}
		pr.Status = model.StatusMerged
		pr.MergedAt = &timeNow

		return nil
	})
	if err != nil{
		return model.PullRequest{}, err
	}
	return pr, nil
}

func (s *pullRequestService) Reassign(ctx context.Context, prID, oldReviewerID string) (model.PullRequest, string, error){
	const pth = "service.pullRequest.Reassign"
	var pr model.PullRequest
	var newReviewer model.User
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		var err error
		pr, err = s.pullRequestRepo.GetByID(ctx, prID)
		if err != nil{
			if errors.Is(err, errs.ErrPRNotFound){
				return errs.ErrPRNotFound
			}
			return fmt.Errorf("getting PR: %w", err)
		}
		exists, err := s.userRepo.IsExistUser(ctx, oldReviewerID)
		if err != nil{
			return fmt.Errorf("getting user: %w", err)
		}
		if !exists{
			return errs.ErrUserNotFound
		}

		if pr.Status == model.StatusMerged{
			return errs.ErrPRMerged
		}

		var oldReviewer model.User
		secReviewer := ""
		for _, u := range pr.AssignedReviewers {
			if u.ID == oldReviewerID {
				oldReviewer = u
			}else{
				secReviewer = u.ID
			}
		}
		if oldReviewer.ID == "" {
			return fmt.Errorf("user %s is not assigned to PR %s", oldReviewerID, prID)
		}

		team, err := s.teamRepo.GetTeamWithMembers(ctx, oldReviewer.TeamName)
		if err != nil {
			return fmt.Errorf("getting team members: %w", err)
		}

		var members []model.User
		for _, u := range team.Members {
			if u.ID != oldReviewerID && u.IsActive != nil && *u.IsActive == true && u.ID != pr.AuthorID && (secReviewer == "" || u.ID != secReviewer){
				members = append(members, u)
			}
		}

		if err := s.pullRequestRepo.UnassignReviewer(ctx, prID, oldReviewerID); err != nil {
			return fmt.Errorf("unassigning old reviewer: %w", err)
		}

		if len(members) != 0 {
			newReviewer := randomUsers(members, 1)[0]
			if err := s.pullRequestRepo.AssignReviewer(ctx, prID, newReviewer.ID); err != nil {
				return fmt.Errorf("assigning new reviewer: %w", err)
			}
		}

		pr, err = s.pullRequestRepo.GetByID(ctx, prID)
		if err != nil {
			return fmt.Errorf("getting PR after reassignment: %w", err)
		}

		return nil
	})
	if err != nil{
		return model.PullRequest{}, "", err
	}
	return pr, newReviewer.ID, nil
}
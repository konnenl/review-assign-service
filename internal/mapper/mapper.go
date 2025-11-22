package mapper

import (
	"github.com/konnen/review-assign-service/internal/dto"
	"github.com/konnen/review-assign-service/internal/model"
)

func UsersToDTO(users []model.User) []dto.UserDTO {
	res := make([]dto.UserDTO, len(users))
	for i, u := range users {
		res[i] = dto.UserDTO{
			ID:       u.ID,
			Name:     u.Name,
			IsActive: u.IsActive,
		}
	}
	return res
}

func UserToWithTeamDTO(u model.User) dto.UserWithTeamDTO {
	return dto.UserWithTeamDTO{
		ID:       u.ID,
		Name:     u.Name,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func TeamToDTO(team model.Team) dto.TeamDTO {
	return dto.TeamDTO{
		Name:    team.Name,
		Members: UsersToDTO(team.Members),
	}
}

func TeamtoModel(team dto.TeamDTO) model.Team {
	members := make([]model.User, len(team.Members))
	for i, u := range team.Members {
		members[i] = model.User{
			ID:       u.ID,
			Name:     u.Name,
			IsActive: u.IsActive,
		}
	}
	return model.Team{
		Name:    team.Name,
		Members: members,
	}
}

func PRtoUserReview(userID string, pullRequests []model.PullRequest) dto.UserReview {
	prs := make([]dto.PullRequestShort, len(pullRequests))
	for i, pr := range pullRequests {
		prs[i] = dto.PullRequestShort{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   string(pr.Status),
		}
	}
	return dto.UserReview{
		ID:           userID,
		PullRequests: prs,
	}
}

func PRtoModel(prDTO dto.PullRequestShort) model.PullRequest {
	return model.PullRequest{
		ID:                prDTO.ID,
		Name:              prDTO.Name,
		AuthorID:          prDTO.AuthorID,
		Status:            model.Status(prDTO.Status),
		AssignedReviewers: nil,
	}
}

func PRtoDTO(pr model.PullRequest) dto.PullRequest {
	reviewerIDs := make([]string, len(pr.AssignedReviewers))
	for i, u := range pr.AssignedReviewers {
		reviewerIDs[i] = u.ID
	}

	prDTO := dto.PullRequest{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: reviewerIDs,
	}

	if pr.Status == model.StatusMerged {
		prDTO.MergedAt = pr.MergedAt
	}
	return prDTO
}

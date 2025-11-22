package dto

type SetIsActiveReq struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResp struct {
	User UserWithTeamDTO `json:"user"`
}

type UserDTO struct {
	ID       string `json:"user_id" validate:"required"`
	Name     string `json:"username" validate:"required"`
	IsActive *bool  `json:"is_active" validate:"is_active"`
}

type UserWithTeamDTO struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive *bool  `json:"is_active"`
}

type UserReview struct {
	ID           string             `json:"user_id" validate:"required"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

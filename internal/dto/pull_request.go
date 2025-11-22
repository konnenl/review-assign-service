package dto

type PullRequestResp struct {
	PullRequest PullRequest `json:"pr"`
}

type PullRequestShort struct {
	ID       string `json:"pull_request_id" validate:"required"`
	Name     string `json:"pull_request_name" validate:"required"`
	AuthorID string `json:"author_id" validate:"required"`
	Status   string `json:"status"`
}

type PullRequest struct {
	ID                string   `json:"pull_request_id" validate:"required"`
	Name              string   `json:"pull_request_name" validate:"required"`
	AuthorID          string   `json:"author_id" validate:"required"`
	Status            string   `json:"status" validate:"required"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

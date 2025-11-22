package dto

import(
	"time"
)

type PullRequestResp struct {
	PullRequest PullRequest `json:"pr"`
	ReplacedBy string `json:"replaced_by,omitempty"`
}

type PullRequestShort struct {
	ID       string `json:"pull_request_id" validate:"required"`
	Name     string `json:"pull_request_name" validate:"required"`
	AuthorID string `json:"author_id" validate:"required"`
	Status   string `json:"status,omitempty"`
}

type PullRequest struct {
	ID                string   `json:"pull_request_id" validate:"required"`
	Name              string   `json:"pull_request_name" validate:"required"`
	AuthorID          string   `json:"author_id" validate:"required"`
	Status            string   `json:"status" validate:"required"`
	AssignedReviewers []string `json:"assigned_reviewers"`
    MergedAt          *time.Time `json:"merged_at,omitempty"`
}

type Merge struct{
	PullRequestID string `json:"pull_request_id" validate:"required"`
}

type Reassign struct{
	PullRequestID string `json:"pull_request_id" validate:"required"`
	OldReviewerID string `json:"old_reviewer_id" validate:"required"`
}
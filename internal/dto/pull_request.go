package dto

type PullRequestShort struct{
	ID string `json:"pull_request_id"`
	Name string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status string `json:"status"`
}
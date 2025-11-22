package model

import(
	"time"
)

type Status string

const(
	StatusMerged  Status = "MERGED"
	StatusOpen Status = "OPEN"
)

type PullRequest struct{
	ID string `db:"id"`
	Name string `db:"name"`
	AuthorID string `db:"author_id"`
	MergedAt *time.Time  `db:"merged_at"`
	Status Status `db:"status"`
	AssignedReviewers []User `db:"-"`
}
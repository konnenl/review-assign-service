package repository

import (
	"fmt"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/konnen/review-assign-service/internal/model"
)

type pullRequestPostgres struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func newPullRequestPostgres(db *sqlx.DB, sq squirrel.StatementBuilderType) *pullRequestPostgres {
	return &pullRequestPostgres{
		db: db,
		sq: sq,
	}
}

func (r *pullRequestPostgres) GetReviews(ctx context.Context, userID string) ([]model.PullRequest, error){
	var prs []model.PullRequest
	query, args, _ := r.sq.
		Select("p.id", "p.name", "p.author_id", "p.merged_at", "status").
		From("pull_requests p").
		Join("reviewers r on p.id = r.pull_request_id").
		Where(squirrel.Eq{"r.reviewer_id": userID}).
		ToSql()
	if err := r.db.SelectContext(ctx, &prs, query, args...); err != nil {
		return nil, fmt.Errorf("executing query GetReviews: %w", err)
	}

	return prs, nil
} 
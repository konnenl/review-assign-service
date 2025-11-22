package repository

import (
	"fmt"
	"context"
	"errors"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/konnen/review-assign-service/internal/model"
	"github.com/konnen/review-assign-service/internal/errs"
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

func (r *pullRequestPostgres) CreatePullRequest(ctx context.Context, pr model.PullRequest) (error){
	query, args, _ := r.sq.
		Insert("pull_requests").
		Columns("id", "name", "author_id", "status").
		Values(pr.ID, pr.Name, pr.AuthorID, string(pr.Status)).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *pullRequestPostgres) GetByID(ctx context.Context, id string) (model.PullRequest, error){
	var pr model.PullRequest
	query, args, _ := r.sq.
		Select("id", "name", "author_id", "merged_at", "status").
		From("pull_request").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	err := r.db.GetContext(ctx, &pr, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.PullRequest{}, errs.ErrPRNotFound
		}
		return model.PullRequest{}, fmt.Errorf("GetByID query: %w", err)
	}
	return pr, nil
}

func (r *pullRequestPostgres) AssignReviewer(ctx context.Context, prID, userID string) (error){
	query, args, _ := r.sq.
		Insert("reviewers").
		Columns("pull_request_id", "reviewer_id").
		Values(prID, userID).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/konnen/review-assign-service/internal/errs"
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

func (r *pullRequestPostgres) GetReviews(ctx context.Context, userID string) ([]model.PullRequest, error) {
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

func (r *pullRequestPostgres) CreatePullRequest(ctx context.Context, pr model.PullRequest) error {
	query, args, _ := r.sq.
		Insert("pull_requests").
		Columns("id", "name", "author_id", "status").
		Values(pr.ID, pr.Name, pr.AuthorID, string(pr.Status)).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *pullRequestPostgres) GetByID(ctx context.Context, id string) (model.PullRequest, error) {
	var pr model.PullRequest
	query, args, _ := r.sq.
		Select("id", "name", "author_id", "merged_at", "status").
		From("pull_requests").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	err := r.db.GetContext(ctx, &pr, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.PullRequest{}, errs.ErrPRNotFound
		}
		return model.PullRequest{}, err
	}

	var reviewers []model.User
	reviewersQuery, reviewersArgs, _ := r.sq.
		Select("u.id", "u.name", "u.is_active", "u.team_name").
		From("reviewers r").
		Join("users u ON r.reviewer_id = u.id").
		Where(squirrel.Eq{"r.pull_request_id": id}).
		ToSql()

	err = r.db.SelectContext(ctx, &reviewers, reviewersQuery, reviewersArgs...)
	if err != nil {
		return model.PullRequest{}, err
	}

	pr.AssignedReviewers = reviewers
	return pr, nil
}

func (r *pullRequestPostgres) AssignReviewer(ctx context.Context, prID, userID string) error {
	query, args, _ := r.sq.
		Insert("reviewers").
		Columns("pull_request_id", "reviewer_id").
		Values(prID, userID).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *pullRequestPostgres) UnassignReviewer(ctx context.Context, prID, userID string) error {
	query, args, _ := r.sq.
		Delete("reviewers").
		Where(squirrel.Eq{"pull_request_id": prID, "reviewer_id": userID}).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *pullRequestPostgres) Merge(ctx context.Context, prID string, timeNow time.Time) error {
	query, args, _ := r.sq.
		Update("pull_requests").
		Set("status", model.StatusMerged).
		Set("merged_at", timeNow).
		Where(squirrel.Eq{"id": prID}).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

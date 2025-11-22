package repository

import (
	"context"
	"errors"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/konnen/review-assign-service/internal/dto"
)

type userPostgres struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func newUserPostgres(db *sqlx.DB, sq squirrel.StatementBuilderType) *userPostgres {
	return &userPostgres{
		db: db,
		sq: sq,
	}
}

func (r *userPostgres) AddUser(ctx context.Context, user dto.User) error {
	query, args, _ := r.sq.Insert("users").
		Columns("id", "name", "is_active").
		Values(user.ID, user.Name, user.IsActive).
		Suffix("ON CONFLICT (id) DO NOTHING").
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *userPostgres) IsExistUser(ctx context.Context, id string) (bool, error) {
	var res int
	query, args, _ := r.sq.Select("1").From("users").Where(squirrel.Eq{"id": id}).ToSql()
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

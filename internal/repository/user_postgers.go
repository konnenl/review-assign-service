package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/konnen/review-assign-service/internal/model"
	"github.com/konnen/review-assign-service/internal/errs"
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

func (r *userPostgres) AddUser(ctx context.Context, user model.User) error {
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
	query, args, _ := r.sq.Select("1").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *userPostgres) SetISActive(ctx context.Context, userID string, isActive bool) (model.User, error) {
	var user model.User
	query, args, _ := r.sq.Update("users").
		Set("is_active", isActive).
		Where(squirrel.Eq{"id": userID}).
		Suffix("RETURNING id, name, is_active, team_name").
		ToSql()
	err := r.db.GetContext(ctx, &user, query, args...)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return user, errs.ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}
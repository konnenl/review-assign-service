package repository

import (
	"context"
	"errors"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/konnen/review-assign-service/internal/dto"
)

type teamPostgres struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func newTeamPostgres(db *sqlx.DB, sq squirrel.StatementBuilderType) *teamPostgres {
	return &teamPostgres{
		db: db,
		sq: sq,
	}
}

func (r *teamPostgres) AddTeam(ctx context.Context, team dto.Team) error {
	query, args, _ := r.sq.Insert("teams").
		Columns("name").
		Values(team.Name).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *teamPostgres) AddMember(ctx context.Context, teamName string, member dto.User) error {
	query, args, _ := r.sq.Insert("users_teams").
		Columns("user_id", "team_name").
		Values(member.ID, teamName).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *teamPostgres) IsTeamExists(ctx context.Context, name string) (bool, error) {
	query, args, _ := r.sq.Select("1").
		From("teams").
		Where(squirrel.Eq{"name": name}).
		ToSql()
	var res int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

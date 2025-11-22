package repository

import (
	"context"
	"errors"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/konnen/review-assign-service/internal/model"
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

func (r *teamPostgres) AddTeam(ctx context.Context, team model.Team) error {
	query, args, _ := r.sq.Insert("teams").
		Columns("name").
		Values(team.Name).
		ToSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *teamPostgres) AddMember(ctx context.Context, teamName string, member model.User) error {
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


func (r *teamPostgres) GetTeamWithMembers(ctx context.Context, teamName string) (model.Team, error){
	var members []model.User
	query, args, _ := r.sq.Select("u.id", "u.name", "u.is_active").
		From("users_teams ut").
		Join("users u ON u.id = ut.user_id").
		Where(squirrel.Eq{"ut.team_name": teamName}).
		ToSql()

	if err := r.db.SelectContext(ctx, &members, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Team{
				Name:    teamName,
				Members: []model.User{},
			}, nil
		}
		return model.Team{}, err
	}

	team := model.Team{
		Name:    teamName,
		Members: members,
	}
	return team, nil
}
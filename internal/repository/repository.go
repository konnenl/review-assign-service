package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Team *teamPostgres
	User *userPostgres
}

func NewRepository(db *sqlx.DB) *Repository {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &Repository{
		Team: newTeamPostgres(db, sq),
		User: newUserPostgres(db, sq),
	}
}

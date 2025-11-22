package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
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

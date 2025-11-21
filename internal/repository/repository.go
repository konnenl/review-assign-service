package repository

import(
	"database/sql"
)

type Repository struct {
	Team *teamPostgres
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Team: newTeamPostgres(db),
	}
}

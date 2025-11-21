package repository

import(
	"database/sql"
)
type teamPostgres struct {
	db *sql.DB
}

func newTeamPostgres(db *sql.DB) *teamPostgres {
	return &teamPostgres{
		db: db,
	}
}

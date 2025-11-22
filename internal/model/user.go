package model

type User struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	IsActive *bool  `db:"is_active"`
	TeamName string `db:"team_name"`
}

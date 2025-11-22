package model

type Team struct {
	Name    string `db:"team_name"`
	Members []User `db:"-"`
}
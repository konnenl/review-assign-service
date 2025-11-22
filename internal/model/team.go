package model

type Team struct {
	Name    string `json:"team_name" validate:"required"`
	Members []User `json:"members" validate:"required,min=1,dive"`
}


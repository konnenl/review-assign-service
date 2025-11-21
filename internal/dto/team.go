package dto

type TeamReq struct{
	Team Team `json:"team" validate:"required"`
}

type Team struct{
	TeamName string `json:"team_name" validate:"required"`
	Members []User `json:"members" validate:"required,min=1,dive"`
}
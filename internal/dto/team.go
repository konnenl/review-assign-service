package dto

type Team struct {
	Name    string `json:"team_name" validate:"required"`
	Members []User `json:"members" validate:"required,min=1,dive"`
}

type TeamResp struct {
	Team Team `json:"team"`
}

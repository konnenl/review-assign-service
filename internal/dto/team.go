package dto

type TeamResp struct {
	Team TeamDTO `json:"team"`
}

type TeamDTO struct {
	Name    string    `json:"team_name"`
	Members []UserDTO `json:"members"`
}

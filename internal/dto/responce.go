package dto

type UserResp struct {
	User UserWithTeamDTO `json:"user"`
}

type TeamResp struct {
	Team TeamDTO `json:"team"`
}

//TODO move from responce.go
type TeamDTO struct {
	Name    string    `json:"team_name"`
	Members []UserDTO `json:"members"`
}

type UserDTO struct{
	ID       string `json:"user_id" validate:"required"`
	Name     string `json:"username" validate:"required"`
	IsActive *bool  `json:"is_active" validate:"is_active"`
}

type UserWithTeamDTO struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive *bool  `json:"is_active"`
}

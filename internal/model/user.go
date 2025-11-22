package model

type User struct {
	ID       string `json:"user_id" validate:"required" db:"id"`
	Name     string `json:"username" validate:"required" db:"name"`
	IsActive *bool  `json:"is_active" validate:"is_active" db:"is_active"`
}

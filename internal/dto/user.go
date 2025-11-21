package dto

type User struct{
	ID string `json:"user_id" validate:"required"`
	Name string `json:"username" validate:"required"`
	IsActive *bool `json:"is_active" validate:"is_active"`
}
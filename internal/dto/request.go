package dto

type SetIsActiveReq struct{
	UserID string `json:"user_id"`
	IsActive bool `json:"is_active"`
}
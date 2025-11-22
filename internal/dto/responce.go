package dto

import (
	"github.com/konnen/review-assign-service/internal/model"
)

type TeamResp struct {
	Team model.Team `json:"team"`
}

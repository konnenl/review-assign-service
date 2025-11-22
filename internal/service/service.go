package service

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type Service struct {
	Team *teamService
}

func NewService(trManager *manager.Manager, teamRepo teamRepository, userRepo userRepository) *Service {
	return &Service{
		Team: newTeamService(trManager, teamRepo, userRepo),
	}
}

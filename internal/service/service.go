package service

type Service struct {
	Team *teamService
}

func NewService(teamRepo teamRepository) *Service {
	return &Service{
		Team: newTeamService(teamRepo),
	}
}

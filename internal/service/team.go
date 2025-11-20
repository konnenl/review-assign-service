package service

type teamRepository interface {
}

type teamService struct {
	teamRepo teamRepository
}

func newTeamService(teamRepo teamRepository) *teamService {
	return &teamService{
		teamRepo: teamRepo,
	}
}

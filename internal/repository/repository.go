package repository

type Repository struct {
	Team *teamPostgres
}

func NewRepository() *Repository {
	return &Repository{
		Team: newTeamPostgres(),
	}
}

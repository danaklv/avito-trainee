package services

import (
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/repository"
)

type TeamService interface {
	GetTeam(team_name string) (*domain.Team, error)
	CreateTeam(team domain.Team) error
}

type teamService struct {
	repo repository.TeamRepository
}

func NewTeamService(r repository.TeamRepository) TeamService {
	return &teamService{repo: r}
}

func (t *teamService) CreateTeam(team domain.Team) error {
	return t.repo.Create(team)
}

func (t *teamService) GetTeam(team_name string) (*domain.Team, error) {
	return t.repo.Get(team_name)

}

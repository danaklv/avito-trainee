package services

import (
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/repository"
)

type StatsService interface {
	GetReviewStats() ([]domain.ReviewerStat, error)
}

type statsService struct {
	repo repository.PullRequestRepository
}

func NewStatsService(r repository.PullRequestRepository) StatsService {
	return &statsService{repo: r}
}

func (s *statsService) GetReviewStats() ([]domain.ReviewerStat, error) {
	return s.repo.GetReviewStats()
}

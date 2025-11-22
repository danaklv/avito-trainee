package services

import (
	"database/sql"
	"errors"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/repository"
)

type UserService interface {
	SetIsActive(userId string, value bool) (*domain.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{userRepo: r}
}

func (s *userService) SetIsActive(userId string, value bool) (*domain.UserResponse, error) {
	user, teamName, err := s.userRepo.GetById(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if user.IsActive == value {
		return nil, domain.ErrAlreadyInState
	}

	if err := s.userRepo.SetIsActive(userId, value); err != nil {
		return nil, err
	}

	user.IsActive = value

	return &domain.UserResponse{
		UserID:   user.ID,
		UserName: user.UserName,
		TeamName: teamName,
		IsActive: user.IsActive,
	}, nil
}

package tests

import (
	"database/sql"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetById(userId string) (*domain.User, string, error) {
	args := m.Called(userId)
	return args.Get(0).(*domain.User), args.String(1), args.Error(2)
}

func (m *MockUserRepository) SetIsActive(userId string, value bool) error {
	args := m.Called(userId, value)
	return args.Error(0)
}


func TestUserService_SetIsActive_OK(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetById", "u1").Return(
		&domain.User{ID: "u1", UserName: "Alice", IsActive: true, TeamID: 1},
		"backend",
		nil,
	)

	mockRepo.On("SetIsActive", "u1", false).Return(nil)

	svc := services.NewUserService(mockRepo)

	_, err := svc.SetIsActive("u1", false)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_SetIsActive_AlreadyState(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetById", "u1").Return(
		&domain.User{ID: "u1", UserName: "Alice", IsActive: false, TeamID: 1},
		"backend",
		nil,
	)

	svc := services.NewUserService(mockRepo)

	_, err := svc.SetIsActive("u1", false)
	assert.ErrorIs(t, err, domain.ErrAlreadyInState)
}

func TestUserService_SetIsActive_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetById", mock.Anything).Return(
		(*domain.User)(nil), "", sql.ErrNoRows,
	)

	svc := services.NewUserService(mockRepo)

	_, err := svc.SetIsActive("404", false)

	assert.ErrorIs(t, err, domain.ErrNotFound)
}

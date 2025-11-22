package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"pr-reviewer/internal/domain"
)

type UserRepository interface {
	SetIsActive(userId string, value bool) error
	GetById(userId string) (*domain.User, string, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) SetIsActive(userId string, value bool) error {

	_, err := u.db.Exec(`UPDATE users SET is_active=$1 WHERE user_id=$2`, value, userId)

	if err != nil {
		return err
	}

	return nil

}

func (u *userRepository) GetById(userId string) (*domain.User, string, error) {

	row := u.db.QueryRow(`
        SELECT u.user_id, u.username, u.is_active, u.team_id, t.team_name
        FROM users u
        JOIN teams t ON u.team_id = t.team_id
        WHERE u.user_id = $1
    `, userId)

	user := &domain.User{}
	var teamName string

	err := row.Scan(&user.ID, &user.UserName, &user.IsActive, &user.TeamID, &teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", sql.ErrNoRows
		}
		return nil, "", fmt.Errorf("select user join teams: %w", err)
	}

	return user, teamName, nil
}

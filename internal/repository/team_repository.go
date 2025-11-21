package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"pr-reviewer/internal/domain"
)

type TeamRepository interface {
	Create(team *domain.Team) error
	Get(team_name string) (*domain.Team, error)
	Exist(team_name string) (bool, error)
}

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(team *domain.Team) error {

	tx, err := r.db.Begin()
	if err != nil {
		return errors.New("tx begin: " + err.Error())
	}

	err = tx.QueryRow(
		"INSERT INTO teams(team_name) VALUES ($1) RETURNING team_id",
		team.TeamName,
	).Scan(&team.ID)

	if err != nil {
		tx.Rollback()
		return errors.New("insert into teams: " + err.Error())
	}

	for _, user := range team.Members {
		log.Println("Inserting:", user.UserName)
		_, err = tx.Exec("INSERT INTO users(user_id, username, is_active, team_id) VALUES ($1, $2, $3, $4)", user.ID, user.UserName, user.IsActive, team.ID)
		if err != nil {
			tx.Rollback()
			return errors.New("insert into users: " + err.Error())
		}

	}

	tx.Commit()
	return nil

}

func (r *teamRepository) Get(team_name string) (*domain.Team, error) {

	var team_id int64

	err := r.db.QueryRow("SELECT team_id FROM teams WHERE team_name=$1", team_name).Scan(&team_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, fmt.Errorf("select from teams: %w", err)
	}

	rows, err := r.db.Query("SELECT user_id, username, is_active FROM users WHERE team_id=$1", team_id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.New("select from users: " + err.Error())
	}

	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var user domain.User

		err := rows.Scan(&user.ID, &user.UserName, &user.IsActive)
		if err != nil {
			return nil, errors.New("scan row: " + err.Error())
		}
		users = append(users, user)
	}

	return &domain.Team{
		ID:       team_id,
		TeamName: team_name,
		Members:  users,
	}, nil

}

func (r *teamRepository) Exist(team_name string) (bool, error) {

	var exist bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`, team_name).Scan(&exist)

	if err != nil {
		err = errors.New("SELECT EXISTS from teams: " + err.Error())
	}
	return exist, err

}

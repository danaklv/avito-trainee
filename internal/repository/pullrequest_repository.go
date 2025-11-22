package repository

import (
	"database/sql"
	"errors"
	"log"
	"pr-reviewer/internal/domain"
	"strings"
)

type PullRequestRepository interface {
	Exists(prID string) (bool, error)
	Create(pr *domain.PullRequest) error
	AssignReviewers(prID string, reviewers []string) error
	GetTeamMembers(teamID int64, exclude string) ([]string, error)
	GetByID(prID string) (*domain.PullRequest, error)
	Merge(prID string, timestamp string) error
	GetReviewers(prID string) ([]string, error)
	ReplaceReviewer(prID, oldID, newID string) error
	FindReplacement(teamID int64, authorID, oldReviewerID string, assigned []string) (string, error)
	GetByReviewer(userID string) ([]domain.PullRequestShort, error)
	GetReviewStats() ([]domain.ReviewerStat, error)
}

type pullRequestRepository struct {
	db *sql.DB
}

func NewPullRequestRepository(db *sql.DB) PullRequestRepository {
	return &pullRequestRepository{db: db}
}

func (r *pullRequestRepository) Exists(prID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id=$1)`, prID).Scan(&exists)
	return exists, err
}

func (r *pullRequestRepository) Create(pr *domain.PullRequest) error {
	_, err := r.db.Exec(`
        INSERT INTO pull_requests (pull_request_id, title, author, status)
        VALUES ($1, $2, $3, $4)
    `, pr.ID, pr.Name, pr.AuthorID, pr.Status)
	return err
}

func (r *pullRequestRepository) AssignReviewers(prID string, reviewers []string) error {
	for _, uid := range reviewers {
		_, err := r.db.Exec(`INSERT INTO reviewers (pull_request_id, user_id) VALUES ($1, $2)`, prID, uid)
		if err != nil {
			return err
		}
	}
	return nil
}

// до 2 участников команды, исключая автора
func (r *pullRequestRepository) GetTeamMembers(teamID int64, exclude string) ([]string, error) {
	rows, err := r.db.Query(`
        SELECT user_id FROM users
        WHERE team_id=$1 AND is_active=true AND user_id != $2
        LIMIT 2
    `, teamID, exclude)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Println("rows close:", cerr)
		}
	}()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *pullRequestRepository) GetByID(prID string) (*domain.PullRequest, error) {
	row := r.db.QueryRow(`
        SELECT pull_request_id, title, author, status, created_at, merged_at
        FROM pull_requests
        WHERE pull_request_id=$1
    `, prID)

	pr := &domain.PullRequest{}
	var mergedAt *string

	err := row.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	pr.MergedAt = mergedAt

	reviewers, err := r.GetReviewers(prID)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return pr, nil
}

func (r *pullRequestRepository) Merge(prID string, timestamp string) error {
	_, err := r.db.Exec(`
        UPDATE pull_requests
        SET status='MERGED', merged_at=$2
        WHERE pull_request_id=$1
    `, prID, timestamp)
	return err
}

func (r *pullRequestRepository) GetReviewers(prID string) ([]string, error) {
	rows, err := r.db.Query(`SELECT user_id FROM reviewers WHERE pull_request_id=$1`, prID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Println("rows close:", cerr)
		}
	}()

	var reviewers []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, id)
	}
	return reviewers, nil
}

func (r *pullRequestRepository) FindReplacement(teamID int64, authorID, oldReviewerID string, assigned []string) (string, error) {

	// исключаем: автора, старого ревьювера, уже назначенных
	query := `
        SELECT user_id
        FROM users
        WHERE team_id = $1
        AND is_active = true
        AND user_id != $2
        AND user_id != $3
        AND user_id NOT IN ($4)
        ORDER BY random()
        LIMIT 1;
    `

	assignedArray := "{" + strings.Join(assigned, ",") + "}"

	var candidate string
	err := r.db.QueryRow(query, teamID, authorID, oldReviewerID, assignedArray).Scan(&candidate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrNoCandidate
		}
		return "", err
	}

	return candidate, nil
}

func (r *pullRequestRepository) ReplaceReviewer(prID, oldID, newID string) error {
	_, err := r.db.Exec(`
        UPDATE reviewers
        SET user_id = $1
        WHERE pull_request_id = $2 AND user_id = $3
    `, newID, prID, oldID)

	return err
}

func (r *pullRequestRepository) GetByReviewer(userID string) ([]domain.PullRequestShort, error) {

	rows, err := r.db.Query(`
        SELECT pr.pull_request_id, pr.title, pr.author, pr.status
        FROM pull_requests pr
        JOIN reviewers r ON pr.pull_request_id = r.pull_request_id
        WHERE r.user_id = $1
    `, userID)

	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Println("rows close:", cerr)
		}
	}()

	var prs []domain.PullRequestShort

	for rows.Next() {
		var pr domain.PullRequestShort
		err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	if len(prs) == 0 {
		return nil, domain.ErrNotFound
	}

	return prs, nil
}

func (r *pullRequestRepository) GetReviewStats() ([]domain.ReviewerStat, error) {
	rows, err := r.db.Query(`
        SELECT user_id, COUNT(*) as cnt
        FROM reviewers
        GROUP BY user_id
        ORDER BY cnt DESC
    `)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Println("rows close:", cerr)
		}
	}()

	var stats []domain.ReviewerStat
	for rows.Next() {
		var s domain.ReviewerStat
		err := rows.Scan(&s.UserID, &s.Count)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}

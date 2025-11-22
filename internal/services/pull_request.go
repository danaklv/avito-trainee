package services

import (
	"errors"
	"log"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/repository"
	"slices"
	"time"
)

type PullRequestService interface {
	Create(pr *domain.PullRequest) (*domain.PullRequest, error)
	Merge(prID string) (*domain.PullRequest, error)
	Reassign(prID, oldReviewerID string) (*domain.PullRequest, string, error)
	GetReview(userID string) ([]domain.PullRequestShort, error)
}

type pullRequestService struct {
	repo  repository.PullRequestRepository
	users repository.UserRepository // get author(user) by id
}

func NewPullRequestService(r repository.PullRequestRepository, ur repository.UserRepository) PullRequestService {
	return &pullRequestService{repo: r, users: ur}
}

func (s *pullRequestService) Create(pr *domain.PullRequest) (*domain.PullRequest, error) {

	exists, err := s.repo.Exists(pr.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrPRExists
	}

	author, _, err := s.users.GetById(pr.AuthorID)
	if err != nil {
		log.Println(err)
		return nil, domain.ErrNotFound
	}

	reviewers, err := s.repo.GetTeamMembers(author.TeamID, author.ID)
	if err != nil {
		return nil, err
	}

	pr.Status = domain.StatusOpen
	pr.AssignedReviewers = reviewers

	if err := s.repo.Create(pr); err != nil {
		log.Println(err)
		return nil, err
	}

	if err := s.repo.AssignReviewers(pr.ID, reviewers); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *pullRequestService) Merge(prID string) (*domain.PullRequest, error) {

	pr, err := s.repo.GetByID(prID)
	if err != nil {
		return nil, err
	}

	//  если уже MERGED
	if pr.Status == domain.StatusMerged {
		return pr, nil
	}

	now := time.Now().UTC().Format(time.RFC3339)

	if err := s.repo.Merge(prID, now); err != nil {
		return nil, err
	}

	pr.Status = domain.StatusMerged
	pr.MergedAt = &now

	return pr, nil
}

func (s *pullRequestService) Reassign(prID, oldReviewerID string) (*domain.PullRequest, string, error) {

	pr, err := s.repo.GetByID(prID)
	if err != nil {
		return nil, "", domain.ErrNotFound
	}

	// если PR уже merged то операция запрещена
	if pr.Status == domain.StatusMerged {
		return nil, "", domain.ErrPRMerged
	}

	// проверяем что старый ревьювер назначен
	found := slices.Contains(pr.AssignedReviewers, oldReviewerID)
	if !found {
		return nil, "", domain.ErrNotAssigned
	}

	
	// прошлый ревьювер
	oldReviewer, _, err := s.users.GetById(oldReviewerID)
	if err != nil {
		return nil, "", domain.ErrNotFound
	}

	// кандидат на замену
	newReviewerID, err := s.repo.FindReplacement(
		oldReviewer.TeamID,
		pr.AuthorID,
		oldReviewerID,
		pr.AssignedReviewers, 
	)
	if err != nil {
		if errors.Is(err, domain.ErrNoCandidate) {
			return nil, "", domain.ErrNoCandidate
		}
		return nil, "", err
	}


	if err := s.repo.ReplaceReviewer(prID, oldReviewerID, newReviewerID); err != nil {
		return nil, "", err
	}


	for idx, r := range pr.AssignedReviewers {
		if r == oldReviewerID {
			pr.AssignedReviewers[idx] = newReviewerID
			break
		}
	}

	return pr, newReviewerID, nil
}

func (s *pullRequestService) GetReview(userID string) ([]domain.PullRequestShort, error) {

	_, _, err := s.users.GetById(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	prs, err := s.repo.GetByReviewer(userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return prs, nil
}

package domain

type Team struct {
	ID       int64  `json:"id"`
	TeamName string `json:"team_name"`
	Members  []User `json:"members"`
}

type User struct {
	ID       string `json:"user_id"`
	UserName string `json:"username"`
	IsActive bool   `json:"is_active"`
	TeamID   int64  `json:"team_id"`
}

type PullRequest struct {
	ID                string   `json:"pull_request_id"`
	Name              string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         string   `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequestShort struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}

type ReviewerStat struct {
	UserID string `json:"user_id"`
	Count  int    `json:"count"`
}

const (
	StatusOpen   = "OPEN"
	StatusMerged = "MERGED"
)

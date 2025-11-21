package domain

type Team struct {
	ID       int64  `json:"id"`
	TeamName string `json:"team_name"`
	Members []User  `json:"members"`
}

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"usermame"`
	IsActive bool   `json:"is_active"`
	TeamID   int64  `json:"team_id"`
}

type PullRequest struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Author User `json:"author"`
	Status string `json:"status"`
}

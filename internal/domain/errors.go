package domain

import "errors"

var (
	ErrNotFound       = errors.New("NOT_FOUND")
	ErrTeamNameTaken  = errors.New("TEAM_EXISTS")
	ErrPRExists       = errors.New("PR_EXISTS")
	ErrAlreadyInState = errors.New("ALREADY_IN_STATE")

	ErrPRMerged    = errors.New("PR_MERGED")
	ErrNotAssigned = errors.New("NOT_ASSIGNED")
	ErrNoCandidate = errors.New("NO_CANDIDATE")
)

type ApiError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func ErrorResponse(code, msg string) ApiError {
	return ApiError{Error: struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{code, msg}}
}

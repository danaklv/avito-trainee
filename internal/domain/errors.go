package domain

import "errors"

var ErrTeamNotFound = errors.New("team not found")
var ErrUserNotFound = errors.New("user not found")
var ErrTeamNameTaken = errors.New("team with this name already exists")


type ApiError struct {
    Error struct {
        Code    string `json:"code"`
        Message string `json:"message"`
    } `json:"error"`
}

func ErrorResponse(code, msg string) ApiError {
    return ApiError{Error: struct{
        Code string `json:"code"`
        Message string `json:"message"`
    }{code, msg}}
}

package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/services"
	"pr-reviewer/internal/utils"
)

type PullRequestHandler struct {
	Service services.PullRequestService
}

func NewPullRequestHandler(s services.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{Service: s}
}

func (h *PullRequestHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.WriteJSON(w, domain.ErrorResponse("METHOD_NOT_ALLOWED", "method not allowed"))
		return
	}

	var body struct {
		ID     string `json:"pull_request_id"`
		Name   string `json:"pull_request_name"`
		Author string `json:"author_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("INVALID_JSON", "invalid json body"))
		return
	}

	if body.ID == "" || body.Name == "" || body.Author == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("VALIDATION_ERROR", "pull_request_id, pull_request_name, author_id required"))
		return
	}

	pr := &domain.PullRequest{
		ID:       body.ID,
		Name:     body.Name,
		AuthorID: body.Author,
	}

	created, err := h.Service.Create(pr)
	if err != nil {

		if errors.Is(err, domain.ErrPRExists) {
			w.WriteHeader(http.StatusConflict)
			utils.WriteJSON(w, domain.ErrorResponse("PR_EXISTS", "PR id already exists"))
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			utils.WriteJSON(w, domain.ErrorResponse("NOT_FOUND", "author or team not found"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJSON(w, domain.ErrorResponse("INTERNAL_ERROR", "failed to create PR"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.WriteJSON(w, map[string]any{
		"pr": created,
	})
}

func (h *PullRequestHandler) Merge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.WriteJSON(w, domain.ErrorResponse("METHOD_NOT_ALLOWED", "method not allowed"))
		return
	}

	var body struct {
		ID string `json:"pull_request_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("INVALID_JSON", "invalid json body"))
		return
	}

	if body.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("VALIDATION_ERROR", "pull_request_id is required"))
		return
	}

	pr, err := h.Service.Merge(body.ID)
	if err != nil {

		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			utils.WriteJSON(w, domain.ErrorResponse("NOT_FOUND", "pull request not found"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJSON(w, domain.ErrorResponse("INTERNAL_ERROR", "failed to merge PR"))
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteJSON(w, map[string]any{
		"pr": pr,
	})
}

func (h *PullRequestHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.WriteJSON(w, domain.ErrorResponse("METHOD_NOT_ALLOWED", "method not allowed"))
		return
	}

	var body struct {
		ID    string `json:"pull_request_id"`
		OldID string `json:"old_user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("INVALID_JSON", "invalid json body"))
		return
	}

	if body.ID == "" || body.OldID == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("VALIDATION_ERROR", "pull_request_id and old_user_id required"))
		return
	}

	pr, replaced, err := h.Service.Reassign(body.ID, body.OldID)
	if err != nil {

		switch {
		case errors.Is(err, domain.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			utils.WriteJSON(w, domain.ErrorResponse("NOT_FOUND", "pr or user not found"))
		case errors.Is(err, domain.ErrPRMerged):
			w.WriteHeader(http.StatusConflict)
			utils.WriteJSON(w, domain.ErrorResponse("PR_MERGED", "cannot reassign on merged PR"))
		case errors.Is(err, domain.ErrNotAssigned):
			w.WriteHeader(http.StatusConflict)
			utils.WriteJSON(w, domain.ErrorResponse("NOT_ASSIGNED", "reviewer is not assigned to this PR"))
		case errors.Is(err, domain.ErrNoCandidate):
			w.WriteHeader(http.StatusConflict)
			utils.WriteJSON(w, domain.ErrorResponse("NO_CANDIDATE", "no active replacement candidate in team"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			utils.WriteJSON(w, domain.ErrorResponse("INTERNAL_ERROR", "failed to reassign reviewer"))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteJSON(w, map[string]any{
		"pr":          pr,
		"replaced_by": replaced,
	})
}

func (h *PullRequestHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.WriteJSON(w, domain.ErrorResponse("METHOD_NOT_ALLOWED", "method not allowed"))
		return
	}

	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("VALIDATION_ERROR", "user_id is required"))
		return
	}

	prs, err := h.Service.GetReview(userID)
	if err != nil {

		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			utils.WriteJSON(w, domain.ErrorResponse("NOT_FOUND", "no pull requests found"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJSON(w, domain.ErrorResponse("INTERNAL_ERROR", "failed to get pull requests"))
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteJSON(w, map[string]any{
		"user_id":       userID,
		"pull_requests": prs,
	})
}

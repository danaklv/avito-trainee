package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/services"
	"pr-reviewer/internal/utils"
)

type UserHandler struct {
	Service services.UserService
}

func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.WriteJSON(w, domain.ErrorResponse("METHOD_NOT_ALLOWED", "method not allowed"))
		return
	}

	var body struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("INVALID_JSON", "invalid json body"))
		return
	}

	if body.UserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("VALIDATION_ERROR", "user_id is required"))
		return
	}

	userDTO, err := h.Service.SetIsActive(body.UserID, body.IsActive)
	if err != nil {

		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			utils.WriteJSON(w, domain.ErrorResponse("NOT_FOUND", "user not found"))
			return
		}

		if errors.Is(err, domain.ErrAlreadyInState) {
			w.WriteHeader(http.StatusConflict)
			utils.WriteJSON(w, domain.ErrorResponse("ALREADY_IN_STATE", "user already in requested state"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJSON(w, domain.ErrorResponse("INTERNAL_ERROR", "failed to update user"))
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteJSON(w, map[string]any{
		"user": userDTO,
	})
}

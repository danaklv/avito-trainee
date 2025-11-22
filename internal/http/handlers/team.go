package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/services"
	"pr-reviewer/internal/utils"
)

type TeamHandler struct {
	Service services.TeamService
}

func NewTeamHandler(s services.TeamService) *TeamHandler {
	return &TeamHandler{Service: s}
}

// CreateTeam handles POST /team/add
func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var team *domain.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("INVALID_JSON", "invalid json body"))
		return
	}

	if team.TeamName == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("VALIDATION_ERROR", "team_name is required"))
		return
	}

	err := h.Service.CreateTeam(team)
	if err != nil {

		if errors.Is(err, domain.ErrTeamNameTaken) {
			w.WriteHeader(http.StatusBadRequest)
			utils.WriteJSON(w, domain.ErrorResponse("TEAM_EXISTS", "team_name already exists"))
			return
		}

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJSON(w, domain.ErrorResponse("INTERNAL_ERROR", "create team failed"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.WriteJSON(w, map[string]*domain.Team{"team": team})
}

// GetTeam handles GET /team/get
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.WriteJSON(w, domain.ErrorResponse("METHOD_NOT_ALLOWED", "method not allowed"))
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSON(w, domain.ErrorResponse("VALIDATION_ERROR", "team_name is required"))
		return
	}

	team, err := h.Service.GetTeam(teamName)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			utils.WriteJSON(w, domain.ErrorResponse("NOT_FOUND", "team not found"))
			return
		}

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJSON(w, domain.ErrorResponse("INTERNAL_ERROR", "failed to get team"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	utils.WriteJSON(w, map[string]any{
		"team": team,
	})
}

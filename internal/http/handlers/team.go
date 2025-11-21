package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/services"
)

type TeamHandler struct {
	Service services.TeamService
}

func NewTeamHandler(s services.TeamService) *TeamHandler {
	return &TeamHandler{Service: s}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var team domain.Team

	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := h.Service.CreateTeam(team)

	if err != nil {
		http.Error(w, "create team failed", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "created"}); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}

}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		http.Error(w, "team_name is required", http.StatusBadRequest)
		return
	}

	team, err := h.Service.GetTeam(teamName)
	if err != nil {

		if errors.Is(err, domain.ErrTeamNotFound) {
			http.Error(w, "team not found", http.StatusNotFound)
			return
		}

		log.Println(err)
		http.Error(w, "failed to get team", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(team); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}

}

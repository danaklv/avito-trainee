package handlers

import (
	"net/http"
	"pr-reviewer/internal/services"
	"pr-reviewer/internal/utils"
)

type StatsHandler struct {
	Service services.StatsService
}

func NewStatsHandler(s services.StatsService) *StatsHandler {
	return &StatsHandler{Service: s}
}

func (h *StatsHandler) GetReviewersStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.Service.GetReviewStats()
	if err != nil {
		http.Error(w, "failed to load stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSON(w, map[string]any{"reviewers": stats})
}

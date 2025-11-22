package api

import (
	"encoding/json"
	"net/http"

	"github.com/wilsonSev/avitoTechService/internal/model"
)

func (h *Handlers) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var newTeam model.Team

	if err := json.NewDecoder(r.Body).Decode(&newTeam); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.CreateTeam(r.Context(), newTeam)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
}

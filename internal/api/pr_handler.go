package api

import (
	"encoding/json"
	"net/http"
	
	"github.com/wilsonSev/avitoTechService/internal/model"
	"github.com/wilsonSev/avitoTechService/internal/services"
	"github.com/wilsonSev/avitoTechService/internal/storage"
)

type PRHandler struct {
	prService *services.PRService
}

func NewPRHandler(prService *services.PRService) *PRHandler {
	return &PRHandler{prService: prService}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {

	var req model.PullRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	pr, err := h.prService.CreatePR(r.Context(), req)

	if err == storage.ErrPRExists {
		writeError(w, http.StatusConflict, "PR_EXISTS", "PR id already exists")
		return
	}
	if err == storage.ErrUserNotFound {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "author not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"pr": pr,
	})
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {

	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	pr, err := h.prService.MergePR(r.Context(), req.PullRequestID)
	if err == storage.ErrPRNotFound {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "PR not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pr": pr,
	})
}

func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {

	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldReviewerID string `json:"old_user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	pr, newReviewer, err := h.prService.ReassignReviewer(
		r.Context(),
		req.PullRequestID,
		req.OldReviewerID,
	)

	if err == storage.ErrPRNotFound {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "PR not found")
		return
	}

	if err != nil {

		switch err.Error() {

		case "cannot reassign on merged PR":
			writeError(w, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
			return

		case "reviewer is not assigned":
			writeError(w, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
			return

		case "no active replacement candidate in team":
			writeError(w, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
			return

		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pr":          pr,
		"replaced_by": newReviewer,
	})
}

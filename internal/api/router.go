package api

import (
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(
	teamHandler *TeamHandler,
	userHandler *UserHandler,
	prHandler *PRHandler,
) *Router {

	mux := http.NewServeMux()

	// TEAM ROUTES
	mux.HandleFunc("POST /team/add", teamHandler.AddTeam)
	mux.HandleFunc("GET /team/get", teamHandler.GetTeam)

	// USER ROUTES
	mux.HandleFunc("POST /users/setIsActive", userHandler.SetIsActive)
	mux.HandleFunc("GET /users/getReview", userHandler.GetUserReviews)

	// PR ROUTES
	mux.HandleFunc("POST /pullRequest/create", prHandler.CreatePR)
	mux.HandleFunc("POST /pullRequest/merge", prHandler.MergePR)
	mux.HandleFunc("POST /pullRequest/reassign", prHandler.ReassignReviewer)

	// healthcheck
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return &Router{mux: mux}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}

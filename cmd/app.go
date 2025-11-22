package main

import (
	"database/sql"
	"log"
	"net/http"
	"pr-reviewer/internal/http/handlers"
	"pr-reviewer/internal/repository"
	"pr-reviewer/internal/services"
	"time"
)

type app struct {
	db   *sql.DB
	port string
}

func NewApp(db *sql.DB, port string) *app {
	return &app{db: db, port: port}
}

func (a *app) Run() {

	// TEAM
	teamRepo := repository.NewTeamRepository(a.db)
	teamService := services.NewTeamService(teamRepo)
	teamHadnler := handlers.NewTeamHandler(teamService)

	// USER
	userRepo := repository.NewUserRepository(a.db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// PULL REQUEST
	pullRequestRepo := repository.NewPullRequestRepository(a.db)
	pullRequestService := services.NewPullRequestService(pullRequestRepo, userRepo)
	pullRequestHandler := handlers.NewPullRequestHandler(pullRequestService)

	// STATS
	statsHandler := handlers.NewStatsHandler(services.NewStatsService(pullRequestRepo))

	http.HandleFunc("/stats/reviewers", statsHandler.GetReviewersStats)

	mux := http.NewServeMux()

	mux.HandleFunc("/users/setIsActive", userHandler.SetIsActive)
	mux.HandleFunc("/users/getReview", pullRequestHandler.GetReview)

	mux.HandleFunc("/team/add", teamHadnler.CreateTeam)
	mux.HandleFunc("/team/get", teamHadnler.GetTeam)

	mux.HandleFunc("/pullRequest/create", pullRequestHandler.CreatePR)
	mux.HandleFunc("/pullRequest/merge", pullRequestHandler.Merge)
	mux.HandleFunc("/pullRequest/reassign", pullRequestHandler.Reassign)

	server := &http.Server{
		Addr:              a.port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("listen and serve on:", a.port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("listen and serve: ", err)
	}

}

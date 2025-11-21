package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pr-reviewer/internal/http/handlers"
	"pr-reviewer/internal/repository"
	"pr-reviewer/internal/services"
	"syscall"
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
	


	mux := http.NewServeMux()

	mux.HandleFunc("/team/add", teamHadnler.CreateTeam)
	mux.HandleFunc("/team/get", teamHadnler.GetTeam)

	server := &http.Server{
		Addr:    a.port,
		Handler: mux,
	}

	go func() {
		log.Println("listen and serve on:", a.port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("listen and serve: ", err)
		}
	}()

	// GRACEFULL SHUTDOWN
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown: ", err)
	}

}

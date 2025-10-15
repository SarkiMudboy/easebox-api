package main

import (
	"log"
	"net/http"

	"github.com/SarkiMudboy/easebox-api/internal/config"
	"github.com/SarkiMudboy/easebox-api/internal/database"
	"github.com/SarkiMudboy/easebox-api/internal/handler"
	"github.com/SarkiMudboy/easebox-api/internal/repository/postgres"
	"github.com/SarkiMudboy/easebox-api/internal/service"
)


func main () {
	cfg := config.Load()

	db, err := database.New(cfg.DB)

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	defer db.Close()

	sessionRepo := postgres.NewSessionRepository(db)
	locationRepo := postgres.NewLocationRepository(db)

	locationService := service.NewLocationService(locationRepo, sessionRepo)

	wsHandler := handler.NewWebSocketHandler(locationService)

	http.HandleFunc("/ws", wsHandler.HandleConnection)
	http.Handle("/", http.FileServer(http.Dir("./web/static")))

	log.Printf("Server starting on %s", cfg.App.ServerAddress)
	if err := http.ListenAndServe(cfg.App.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
package main

import (
	"log"
	"net/http"

	"github.com/SarkiMudboy/easebox-api/internal/config"
	"github.com/SarkiMudboy/easebox-api/internal/database"
)


func main () {
	cfg := config.Load()

	db, err := database.New(cfg.DB)

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	defer db.Close()


	log.Printf("Server starting on %s", cfg.App.ServerAddress)
	if err := http.ListenAndServe(cfg.App.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/michel991/go-api-tech-challenge/internal/api"
	"github.com/michel991/go-api-tech-challenge/internal/config"
	"github.com/michel991/go-api-tech-challenge/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Loaded configuration: Database{Name: %s, User: %s, Password: ****, Host: %s, Port: %s}, HTTPPort: %s",
		cfg.Database.Name, cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.HTTPPort)

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize the chi router
	r := chi.NewRouter()

	// Setup routes
	api.SetupRoutes(r, db)

	// Start the HTTP server
	log.Printf("Starting server on port %s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

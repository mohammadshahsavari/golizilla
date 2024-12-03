package main

import (
	"log"

	"golizilla/config"
	"golizilla/handler/middleware"
	"golizilla/persistence/database"
	"golizilla/route"

	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	// Setup database with retry logic
	db, err := database.SetupDatabaseWithRetry(cfg, 5)
	if err != nil {
		log.Printf("Failed to set up the database after multiple attempts: %v", err)
		return
	}

	// Type assertion to *gorm.DB
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		log.Printf("Failed to assert database to *gorm.DB")
		return
	}

	// Start API
	middleware.InitSessionStore(cfg)
	route.RunServer(cfg, gormDB)

}

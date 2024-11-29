package main

import (
	"log"

	"golizilla/config"
	"golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/persistence/database"
	"golizilla/route"
	"golizilla/service"

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

	// Initialize repositories
	userRepo := repository.NewUserRepository(gormDB)

	// Initialize services
	emailService := service.NewEmailService(cfg)
	userService := service.NewUserService(userRepo, emailService)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, emailService, cfg)

	// Start API
	route.RunServer(cfg, userHandler)

}

package main

import (
	"fmt"
	"log"

	"golizilla/config"
	"golizilla/handler/middleware"
	"golizilla/persistence/database"
	"golizilla/persistence/logger"
	"golizilla/route"

	"go.uber.org/zap/zapcore"
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

	// Initialize the singleton logger
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.MongoDbUsername, cfg.MongoDbPassword, cfg.MongoDbHost, cfg.MongoDbPort)
	if err := logger.Initialize(mongoURI, "logsdb", "logs", zapcore.InfoLevel); err != nil {
		log.Println("Failed to initialize logger:", err)
		return
	}

	// Start API
	middleware.InitSessionStore(cfg)
	route.RunServer(cfg, gormDB)

}

package database

import (
	"fmt"
	"golizilla/config"
	models "golizilla/domain/model"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// SetupDB initializes and returns a database connection
func setupDB(cfg *config.Config) (*gorm.DB, error) {
	// Retrieve the database URL from environment variables
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Run migrations for the Ads, Filters, and Users models
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// TODO:  Check if Super Admin exists, and create one if not

	return db, nil
}

func SetupDatabaseWithRetry(cfg *config.Config, retries int) (interface{}, error) {
	var db interface{}
	var err error
	for retries > 0 {
		db, err = setupDB(cfg)
		if err == nil {
			return db, nil
		}
		retries--
		log.Printf("Error setting up the database, retrying... (%d retries left)", retries)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

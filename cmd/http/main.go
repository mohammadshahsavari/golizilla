package main

import (
	"fmt"
	"log"
	"time"

	"golizilla/config"
	"golizilla/persistence/database"
	"golizilla/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
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

	// Initialize Fiber app with middleware
	app := fiber.New()

	// Add middleware for logging, panic recovery, and CORS
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Add additional middleware for rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	// Grouping routes for users
	userGroup := app.Group("/users")
	route.SetupUserRoutes(userGroup, gormDB)

	// Serve Swagger UI if needed, restricted to non-production environments
	if cfg.Env != "production" {
		// Serve static files from the docs directory if needed, with rate limiting
		app.Use("/docs", limiter.New(limiter.Config{
			Max:        10,
			Expiration: 30 * time.Second,
		}))
		app.Static("/docs", "./docs")

		app.Get("/swagger/*", swagger.New(swagger.Config{
			URL: "/docs/swagger.json", // The URL where swagger.json is located
		}))
	}

	// Start the server
	host := cfg.Host
	port := cfg.Port
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", host, port)))
}

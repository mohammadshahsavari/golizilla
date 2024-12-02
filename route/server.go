package route

import (
	"fmt"
	"log"
	"time"

	"golizilla/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func RunServer(cfg *config.Config, database *gorm.DB) {
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
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as the key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many requests. Please try again later.",
			})
		},
	}))

	// Setup routes
	SetupUserRoutes(app, database, cfg)

	setupQuestionnariRoutes(app, database, cfg)

	// Start the server
	host := cfg.Host
	port := cfg.Port
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", host, port)))
}

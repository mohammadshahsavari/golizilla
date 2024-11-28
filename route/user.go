package route

import (
	"golizilla/config"
	"golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/handler/middleware"
	"golizilla/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// Grouping routes for users
	userGroup := app.Group("/users")
	// Create repository
	userRepo := repository.NewUserRepository(db)

	// Create services
	userService := service.NewUserService(userRepo)
	emailService := service.NewEmailService(cfg)

	// Create handler
	userHandler := handler.NewUserHandler(userService, emailService, cfg)

	// Public routes
	userGroup.Post("/signup", userHandler.CreateUser)
	userGroup.Post("/verify-signup", userHandler.VerifySignup)

	rateLimitConfig := limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many login attempts. Please try again later.",
			})
		},
	}

	userGroup.Post("/login", limiter.New(rateLimitConfig), userHandler.Login)
	userGroup.Post("/verify-login", userHandler.VerifyLogin)

	// Protected routes
	userGroup.Use(middleware.JWTMiddleware(cfg))
	userGroup.Get("/profile", userHandler.GetProfile)
	userGroup.Post("/enable-2fa", userHandler.Enable2FA)
	userGroup.Post("/disable-2fa", userHandler.Disable2FA)
	userGroup.Post("/logout", userHandler.Logout)
}

package route

import (
	"golizilla/config"
	"golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/handler/middleware"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
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
	userGroup.Post("/login", userHandler.Login)
	userGroup.Post("/verify-login", userHandler.VerifyLogin)

	// Protected routes
	userGroup.Use(middleware.JWTMiddleware(cfg))
	userGroup.Get("/profile", userHandler.GetProfile)
}

package route

import (
	"golizilla/config"
	"golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// Grouping routes for users
	userGroup := app.Group("/users")
	// Create repository
	userRepo := repository.NewUserRepository(db)

	// Create service with the repository
	userService := service.NewUserService(userRepo)

	// Create email service
	emailService := service.NewEmailService(cfg)

	// Create handler with the service and email service
	userHandler := handler.NewUserHandler(userService, emailService)

	// Define routes
	userGroup.Post("/signup", userHandler.CreateUser)
	userGroup.Post("/verify-signup", userHandler.VerifySignup)
	userGroup.Get("/:id", userHandler.GetUserByID)
}

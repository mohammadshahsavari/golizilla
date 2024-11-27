package route

import (
	"golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, db *gorm.DB) {
	// Grouping routes for users
	userGroup := app.Group("/users")
	// Create repository
	userRepo := repository.NewUserRepository(db)

	// Create service with the repository
	userService := service.NewUserService(userRepo)

	// Create handler with the service
	userHandler := handler.NewUserHandler(userService)

	// Define routes
	userGroup.Post("/", userHandler.CreateUser)
	userGroup.Get("/:id", userHandler.GetUserByID)
}

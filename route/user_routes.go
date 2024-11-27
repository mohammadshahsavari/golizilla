package route

import (
	"golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(app fiber.Router, db *gorm.DB) {
	// Create repository
	userRepo := repository.NewUserRepository(db)

	// Create service with the repository
	userService := service.NewUserService(userRepo)

	// Create handler with the service
	userHandler := handler.NewUserHandler(userService)

	// Define routes
	app.Post("/", userHandler.CreateUser)
	app.Get("/:id", userHandler.GetUserByID)
}

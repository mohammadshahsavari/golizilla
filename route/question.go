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

func SetupQuestionRoutes(app *fiber.App, database *gorm.DB, cfg *config.Config) {
	// Create a group for user routes
	questionGroup := app.Group("/questions")

	// Initialize repositories
	questionRepo := repository.NewQuestionRepository(database)

	// Initialize services
	questionService := service.NewQuestionService(questionRepo)

	// Initialize handlers
	questionHandler := handler.NewQuestionHandler(questionService)

	// Initialize the JWT middleware with the config
	jwtMiddleware := middleware.AuthMiddleware(cfg)

	// Protected routes
	questionGroup.Post("/create", jwtMiddleware, questionHandler.Create)
	questionGroup.Put("/update/:id", jwtMiddleware, questionHandler.Update)
	questionGroup.Get("/:id", jwtMiddleware, questionHandler.GetByID)
	questionGroup.Delete("/:id", jwtMiddleware, questionHandler.Delete)
}
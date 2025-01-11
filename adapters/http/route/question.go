package route

import (
	"golizilla/adapters/http/handler"
	"golizilla/adapters/http/handler/middleware"
	"golizilla/config"
	"golizilla/core/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupQuestionRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, questionService service.IQuestionService) {
	// Create a group for user routes
	questionGroup := app.Group("/question")

	// Initialize handlers
	questionHandler := handler.NewQuestionHandler(questionService)

	// Initialize the JWT middleware with the config
	questionGroup.Use(middleware.AuthMiddleware(cfg))
	questionGroup.Use(middleware.ContextMiddleware())

	// Protected routes
	questionGroup.Post("/create", questionHandler.Create)

	questionGroup.Put("/update/:id", questionHandler.Update)

	questionGroup.Get("/:id", questionHandler.GetByID)

	questionGroup.Delete("/:id", questionHandler.Delete)
}

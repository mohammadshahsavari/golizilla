package route

import (
	"golizilla/config"
	"golizilla/handler"
	"golizilla/handler/middleware"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupQuestionRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, questionService service.IQuestionService) {
	// Create a group for user routes
	questionGroup := app.Group("/question")

	// Initialize handlers
	questionHandler := handler.NewQuestionHandler(questionService)

	// Initialize the JWT middleware with the config
	authMiddleware := middleware.AuthMiddleware(cfg)

	// Protected routes
	questionGroup.Post("/create", authMiddleware, questionHandler.Create)

	questionGroup.Put("/update/:id", authMiddleware, questionHandler.Update)

	questionGroup.Get("/:id", authMiddleware, questionHandler.GetByID)

	questionGroup.Delete("/:id", authMiddleware, questionHandler.Delete)
}

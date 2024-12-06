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
	questionGroup.Post("/create",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionHandler.Create)

	questionGroup.Put("/update/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionHandler.Update)

	questionGroup.Get("/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionHandler.GetByID)

	questionGroup.Delete("/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionHandler.Delete)
}

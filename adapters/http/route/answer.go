package route

import (
	"golizilla/adapters/http/handler"
	"golizilla/adapters/http/handler/middleware"
	"golizilla/config"
	"golizilla/core/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAnswerRoutes(
	app *fiber.App,
	db *gorm.DB,
	cfg *config.Config,
	answerService service.IAnswerService,
) {
	// Create a group for user routes
	answerGroup := app.Group("/answer")

	// Initialize handlers
	answerHandler := handler.NewAnswerHandler(answerService)

	// Initialize the JWT middleware with the config
	answerGroup.Use(middleware.AuthMiddleware(cfg))
	answerGroup.Use(middleware.ContextMiddleware())

	// Protected routes
	answerGroup.Post("/create", answerHandler.Create)
	answerGroup.Put("/update/:id", answerHandler.Update)
	answerGroup.Get("/:id", answerHandler.GetByID)
	answerGroup.Delete("/:id", answerHandler.Delete)
}

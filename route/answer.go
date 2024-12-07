package route

import (
	"golizilla/config"
	"golizilla/handler"
	"golizilla/handler/middleware"
	"golizilla/service"

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
	authMiddleware := middleware.AuthMiddleware(cfg)

	// Protected routes
	answerGroup.Post("/create", authMiddleware, answerHandler.Create)

	answerGroup.Put("/update/:id", authMiddleware, answerHandler.Update)
	answerGroup.Get("/:id", authMiddleware, answerHandler.GetByID)

	answerGroup.Delete("/:id", authMiddleware, answerHandler.Delete)
}

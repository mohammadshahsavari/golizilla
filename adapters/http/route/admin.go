package route

import (
	"golizilla/adapters/http/handler"
	"golizilla/adapters/http/handler/middleware"
	"golizilla/config"
	"golizilla/core/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAdminRoutes(
	app *fiber.App,
	db *gorm.DB,
	cfg *config.Config,
	adminService service.IAdminService,
) {
	// Create a group for user routes
	adminGroup := app.Group("/admin")

	// Initialize handlers
	adminHandler := handler.NewAdminHandler(adminService)

	// Initialize the JWT middleware with the config
	authMiddleware := middleware.AuthMiddleware(cfg)

	// Protected routes
	adminGroup.Get("/users", authMiddleware, adminHandler.GetAllUsers)
	adminGroup.Get("/questions", authMiddleware, adminHandler.GetAllQuestions)
	adminGroup.Get("/questionnaires", authMiddleware, adminHandler.GetAllQuestionnaires)
	adminGroup.Get("/roles", authMiddleware, adminHandler.GetAllRoles)
	adminGroup.Get("/users/:userID/questionnaires/:questionnaireID", authMiddleware, adminHandler.GetAnswersByUserIDAndQuestionnaireID)
}

package route

import (
	"golizilla/config"
	"golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/handler/middleware"
	"golizilla/persistence/database"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupQuestionRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// Create a group for user routes
	questionGroup := app.Group("/questions")

	// Initialize repositories
	questionRepo := repository.NewQuestionRepository(db)

	// Initialize services
	questionService := service.NewQuestionService(questionRepo)

	// Initialize handlers
	questionHandler := handler.NewQuestionHandler(questionService)

	// Initialize the JWT middleware with the config
	jwtMiddleware := middleware.AuthMiddleware(cfg)

	// Protected routes
	questionGroup.Post("/create", middleware.SetTransaction(database.NewGormCommitter(db)), jwtMiddleware, questionHandler.Create)
	questionGroup.Put("/update/:id", middleware.SetTransaction(database.NewGormCommitter(db)), jwtMiddleware, questionHandler.Update)
	questionGroup.Get("/:id", middleware.SetTransaction(database.NewGormCommitter(db)), jwtMiddleware, questionHandler.GetByID)
	questionGroup.Delete("/:id", middleware.SetTransaction(database.NewGormCommitter(db)), jwtMiddleware, questionHandler.Delete)
}

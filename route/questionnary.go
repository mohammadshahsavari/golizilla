package route

import (
	"golizilla/config"
	respository "golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/handler/middleware"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func setupQuestionnariRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	questionnaryGroup := app.Group("/questionnari")

	questionnaryRepo := respository.NewQuestionnaireRepository(db)

	questionnaryService := service.NewQuestionnaireService(questionnaryRepo)

	questionnaryHandler := handler.NewQuestionnaryHandler(questionnaryService)

	authMiddleware := middleware.AuthMiddleware(cfg)

	questionnaryGroup.Post("/", authMiddleware, questionnaryHandler.Create)

	questionnaryGroup.Get("/:id", authMiddleware, questionnaryHandler.GetById)
	questionnaryGroup.Get("/ownerId/:id", authMiddleware, questionnaryHandler.GetByOwnerId)
	questionnaryGroup.Post("/update", authMiddleware, questionnaryHandler.Update)
	questionnaryGroup.Delete("/:id", authMiddleware, questionnaryHandler.Delete)

}

package route

import (
	"golizilla/config"
	respository "golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/handler/middleware"
	"golizilla/service"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func setupQuestionnariRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	questionnaryGroup := app.Group("/questionnari")

	questionnaryRepo := respository.NewQuestionnaireRepository(db)

	questionnaryService := service.NewQuestionnaireService(questionnaryRepo)

	questionnariHandler := handler.NewQuestionnaryHandler(questionnaryService)

	authMiddleware := middleware.AuthMiddleware(cfg)
	headerAuthMiddleware := middleware.HeaderAuthMiddleware(cfg)

	questionnaryGroup.Post("/", authMiddleware, questionnariHandler.Create)

	questionnaryGroup.Get("/:id", authMiddleware, questionnariHandler.GetById)
	questionnaryGroup.Get("/ownerId/:id", authMiddleware, questionnariHandler.GetByOwnerId)
	questionnaryGroup.Post("/update", authMiddleware, questionnariHandler.Update)
	questionnaryGroup.Delete("/:id", authMiddleware, questionnariHandler.Delete)
	questionnaryGroup.Get("/GetResults/:id", headerAuthMiddleware, websocket.New(questionnariHandler.GetResults))
}

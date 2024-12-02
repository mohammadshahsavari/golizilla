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
	questionnariGroup := app.Group("/questionnari")

	questionnariRepo := respository.NewQuestionnaireRepository(db)

	questionnariService := service.NewQuestionnaireService(questionnariRepo)

	questionnariHandler := handler.NewQuestionnariHandler(questionnariService)

	authMiddleware := middleware.AuthMiddleware(cfg)
	headerAuthMiddleware := middleware.HeaderAuthMiddleware(cfg)

	questionnariGroup.Post("/", authMiddleware, questionnariHandler.Create)

	questionnariGroup.Get("/:id", authMiddleware, questionnariHandler.GetById)
	questionnariGroup.Get("/ownerId/:id", authMiddleware, questionnariHandler.GetByOwnerId)
	questionnariGroup.Post("/update", authMiddleware, questionnariHandler.Update)
	questionnariGroup.Delete("/:id", authMiddleware, questionnariHandler.Delete)
	questionnariGroup.Get("/GetResults/:id", headerAuthMiddleware, websocket.New(questionnariHandler.GetResults))
}

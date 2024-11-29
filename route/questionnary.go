package route

import (
	respository "golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func setupQuestionnariRoutes(app *fiber.App, db *gorm.DB) {
	questionnaryGroup := app.Group("/questionnari")

	questionnaryRepo := respository.NewQuestionnaireRepository(db)

	questionnaryService := service.NewQuestionnaireService(questionnaryRepo)

	questionnaryHandler := handler.NewQuestionnaryHandler(questionnaryService)

	questionnaryGroup.Post("/", questionnaryHandler.Create)

	questionnaryGroup.Get("/:id", questionnaryHandler.GetById)
	questionnaryGroup.Get("/ownerId/:id", questionnaryHandler.GetByOwnerId)
	questionnaryGroup.Post("/update", questionnaryHandler.Update)
	questionnaryGroup.Delete("/:id", questionnaryHandler.Delete)

}

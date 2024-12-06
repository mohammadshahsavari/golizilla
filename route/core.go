package route

import (
	"golizilla/config"
	"golizilla/handler"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupCoreRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	coreGroup := app.Group("/core")

	coreService := service.NewCoreService()

	coreHandler := handler.NewCoreHandler(coreService)

	coreGroup.Post("/start/:questionnaire_id", coreHandler.StartHandler)
	coreGroup.Post("/submit/:question_id", coreHandler.SubmitHandler)
	coreGroup.Post("/back", coreHandler.BackHandler)
	coreGroup.Post("/next", coreHandler.NextHandler)
	coreGroup.Post("/end", coreHandler.EndHandler)
}

package route

import (
	"golizilla/config"
	"golizilla/handler"
	"golizilla/handler/middleware"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupCoreRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config,
	coreService service.ICoreService) { // pass the service as a param or create inside
	coreGroup := app.Group("/core")

	coreHandler := handler.NewCoreHandler(coreService)

	// Add authentication middleware if needed
	coreGroup.Use(middleware.AuthMiddleware(cfg))
	coreGroup.Use(middleware.ContextMiddleware())

	coreGroup.Get("/start/:questionnaire_id", coreHandler.StartHandler)
	coreGroup.Post("/submit", coreHandler.SubmitHandler)
	coreGroup.Post("/next", coreHandler.NextHandler)
	coreGroup.Post("/back", coreHandler.BackHandler)
	coreGroup.Post("/end", coreHandler.EndHandler)
}

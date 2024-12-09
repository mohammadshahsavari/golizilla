package route

import (
	"golizilla/adapters/http/handler"
	"golizilla/adapters/http/handler/middleware"
	"golizilla/config"
	"golizilla/core/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupCoreRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config,
	coreService service.ICoreService, roleService service.IRoleService, questionnaireService service.IQuestionnaireService) { // pass the service as a param or create inside
	coreGroup := app.Group("/core")

	coreHandler := handler.NewCoreHandler(coreService, roleService, questionnaireService)

	// Add authentication middleware if needed
	coreGroup.Use(middleware.AuthMiddleware(cfg))
	coreGroup.Use(middleware.ContextMiddleware())

	coreGroup.Get("/start/:questionnaire_id", coreHandler.StartHandler)
	coreGroup.Post("/submit", coreHandler.SubmitHandler)
	coreGroup.Post("/next", coreHandler.NextHandler)
	coreGroup.Post("/back", coreHandler.BackHandler)
	coreGroup.Post("/end", coreHandler.EndHandler)
}

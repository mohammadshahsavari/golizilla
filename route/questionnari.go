package route

import (
	"golizilla/config"
	"golizilla/handler"
	"golizilla/handler/middleware"
	privilegeconstants "golizilla/internal/privilegeConstants"
	"golizilla/service"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupQuestionnaireRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, questionnaireService service.IQuestionnaireService, authorizationService service.IAuthorizationService) {
	questionnaireGroup := app.Group("/questionnaire")

	questionnaireHandler := handler.NewQuestionnaireHandler(questionnaireService)

	authMiddleware := middleware.AuthMiddleware(cfg)
	authorizationMiddleware := middleware.AuthorizationMiddleware(authorizationService)
	headerAuthMiddleware := middleware.HeaderAuthMiddleware(cfg)

	questionnaireGroup.Post("/",
		authMiddleware, authorizationMiddleware(privilegeconstants.CreateQuestionnaire), questionnaireHandler.Create)

	questionnaireGroup.Get("/:id",
		authMiddleware, questionnaireHandler.GetById)

	questionnaireGroup.Get("/ownerId/:id",
		authMiddleware, questionnaireHandler.GetByOwnerId)

	questionnaireGroup.Put("/update/:id",
		authMiddleware, questionnaireHandler.Update)

	questionnaireGroup.Delete("/:id",
		authMiddleware, questionnaireHandler.Delete)

	questionnaireGroup.Get("/GetResults/:id",
		headerAuthMiddleware, websocket.New(questionnaireHandler.GetResults))
}

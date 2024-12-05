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
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, authorizationMiddleware(privilegeconstants.CreateQuestionnaire), questionnaireHandler.Create)

	questionnaireGroup.Get("/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionnaireHandler.GetById)

	questionnaireGroup.Get("/ownerId/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionnaireHandler.GetByOwnerId)

	questionnaireGroup.Post("/update",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionnaireHandler.Update)

	questionnaireGroup.Delete("/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionnaireHandler.Delete)

	questionnaireGroup.Get("/GetResults/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		headerAuthMiddleware, websocket.New(questionnaireHandler.GetResults))
}

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

func SetupQuestionnariRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, questionnariService service.IQuestionnaireService, authorizationService service.IAuthorizationService) {
	questionnariGroup := app.Group("/questionnari")

	questionnariHandler := handler.NewQuestionnariHandler(questionnariService)

	authMiddleware := middleware.AuthMiddleware(cfg)
	authorizationMiddleware := middleware.AuthorizationMiddleware(authorizationService)
	headerAuthMiddleware := middleware.HeaderAuthMiddleware(cfg)

	questionnariGroup.Post("/",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, authorizationMiddleware(privilegeconstants.CreateQuestionnari), questionnariHandler.Create)

	questionnariGroup.Get("/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionnariHandler.GetById)

	questionnariGroup.Get("/ownerId/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionnariHandler.GetByOwnerId)

	questionnariGroup.Post("/update",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionnariHandler.Update)

	questionnariGroup.Delete("/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		authMiddleware, questionnariHandler.Delete)

	questionnariGroup.Get("/GetResults/:id",
		// middleware.SetTransaction(database.NewGormCommitter(db)),
		headerAuthMiddleware, websocket.New(questionnariHandler.GetResults))
}

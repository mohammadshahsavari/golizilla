package route

import (
	"golizilla/adapters/http/handler"
	"golizilla/adapters/http/handler/middleware"
	"golizilla/config"
	"golizilla/core/service"
	privilegeconstants "golizilla/internal/privilege"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupQuestionnaireRoutes(
	app *fiber.App,
	db *gorm.DB,
	cfg *config.Config,
	questionnaireService service.IQuestionnaireService,
	authorizationService service.IAuthorizationService,
	roleService service.IRoleService,
	userService service.IUserService,
	questionService service.IQuestionService) {
	questionnaireGroup := app.Group("/questionnaire")

	questionnaireHandler := handler.NewQuestionnaireHandler(questionnaireService, roleService, userService, questionService)

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

	questionnaireGroup.Post("/GiveAcess/:id", authMiddleware, questionnaireHandler.GiveAcess)

	questionnaireGroup.Post("/DeleteAcess/:id", authMiddleware, questionnaireHandler.DeleteAcess)
}

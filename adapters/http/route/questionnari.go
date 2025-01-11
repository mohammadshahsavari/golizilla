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

	headerAuthMiddleware := middleware.HeaderAuthMiddleware(cfg)
	questionnaireGroup.Get("/GetResults/:id",
		headerAuthMiddleware, websocket.New(questionnaireHandler.GetResults))

	questionnaireGroup.Use(middleware.AuthMiddleware(cfg))
	authorizationMiddleware := middleware.AuthorizationMiddleware(authorizationService)

	questionnaireGroup.Use(middleware.ContextMiddleware())

	questionnaireGroup.Post("/",
		authorizationMiddleware(privilegeconstants.CreateQuestionnaire), questionnaireHandler.Create)

	questionnaireGroup.Get("/:id",
		questionnaireHandler.GetById)

	questionnaireGroup.Get("/ownerId/:id",
		questionnaireHandler.GetByOwnerId)

	questionnaireGroup.Put("/update/:id",
		questionnaireHandler.Update)

	questionnaireGroup.Delete("/:id",
		questionnaireHandler.Delete)

	questionnaireGroup.Post("/GiveAcess/:id", questionnaireHandler.GiveAcess)

	questionnaireGroup.Post("/DeleteAcess/:id", questionnaireHandler.DeleteAcess)
}

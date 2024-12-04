package route

import (
	"golizilla/config"
	"golizilla/domain/model"
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
	roleRepo := respository.NewRoleRepository(db)
	userRepo := respository.NewUserRepository(db)
	rolePrivilegeRepo := respository.NewRolePrivilegeRepository(db)

	questionnariService := service.NewQuestionnaireService(questionnariRepo)
	roleService := service.NewRoleService(roleRepo, userRepo, rolePrivilegeRepo)
	authorizationsService := service.NewAuthorizationService(roleService)

	questionnariHandler := handler.NewQuestionnariHandler(questionnariService)

	authMiddleware := middleware.AuthMiddleware(cfg)
	authorizationMiddleware := middleware.AuthorizationMiddleware(authorizationsService)
	headerAuthMiddleware := middleware.HeaderAuthMiddleware(cfg)

	questionnariGroup.Post("/", authMiddleware, authorizationMiddleware(model.CreateQuestionnari), questionnariHandler.Create)

	questionnariGroup.Get("/:id", authMiddleware, questionnariHandler.GetById)
	questionnariGroup.Get("/ownerId/:id", authMiddleware, questionnariHandler.GetByOwnerId)
	questionnariGroup.Post("/update", authMiddleware, questionnariHandler.Update)
	questionnariGroup.Delete("/:id", authMiddleware, questionnariHandler.Delete)
	questionnariGroup.Get("/GetResults/:id", headerAuthMiddleware, websocket.New(questionnariHandler.GetResults))
}

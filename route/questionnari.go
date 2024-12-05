package route

import (
	"golizilla/config"
	respository "golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/handler/middleware"
	privilegeconstants "golizilla/internal/privilegeConstants"
	"golizilla/persistence/database"
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

	questionnariGroup.Post("/", middleware.SetTransaction(database.NewGormCommitter(db)), authMiddleware, authorizationMiddleware(privilegeconstants.CreateQuestionnari), questionnariHandler.Create)
	questionnariGroup.Get("/:id", middleware.SetTransaction(database.NewGormCommitter(db)), authMiddleware, questionnariHandler.GetById)
	questionnariGroup.Get("/ownerId/:id", middleware.SetTransaction(database.NewGormCommitter(db)), authMiddleware, questionnariHandler.GetByOwnerId)
	questionnariGroup.Post("/update", middleware.SetTransaction(database.NewGormCommitter(db)), authMiddleware, questionnariHandler.Update)
	questionnariGroup.Delete("/:id", middleware.SetTransaction(database.NewGormCommitter(db)), authMiddleware, questionnariHandler.Delete)
	questionnariGroup.Get("/GetResults/:id", middleware.SetTransaction(database.NewGormCommitter(db)), headerAuthMiddleware, websocket.New(questionnariHandler.GetResults))
}

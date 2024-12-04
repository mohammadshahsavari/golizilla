package route

import (
	"golizilla/config"
	"golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/handler/middleware"
	"golizilla/persistence/database"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// Create a group for user routes
	userGroup := app.Group("/users")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	rolePrivilege := repository.NewRolePrivilegeRepository(db)

	// Initialize services
	emailService := service.NewEmailService(cfg)
	userService := service.NewUserService(userRepo, emailService)
	roleService := service.NewRoleService(roleRepo, userRepo, rolePrivilege)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, emailService, roleService, cfg)

	// Public routes
	userGroup.Post("/register", middleware.SetTransaction(database.NewGormCommitter(db)), userHandler.CreateUser)
	userGroup.Post("/login", middleware.SetTransaction(database.NewGormCommitter(db)), userHandler.Login)
	userGroup.Post("/verify-email", middleware.SetTransaction(database.NewGormCommitter(db)), userHandler.VerifySignup)
	userGroup.Post("/verify-login", middleware.SetTransaction(database.NewGormCommitter(db)), userHandler.VerifyLogin)

	// Initialize the JWT middleware with the config
	jwtMiddleware := middleware.AuthMiddleware(cfg)

	// Protected routes
	userGroup.Get("/profile", middleware.SetTransaction(database.NewGormCommitter(db)), jwtMiddleware, userHandler.GetProfile)
	userGroup.Put("/profile/update", middleware.SetTransaction(database.NewGormCommitter(db)), jwtMiddleware, userHandler.UpdateProfile)
	userGroup.Get("/profile/notifications", middleware.SetTransaction(database.NewGormCommitter(db)), jwtMiddleware, userHandler.GetNotificationListList)
	userGroup.Post("/enable-2fa", middleware.SetTransaction(database.NewGormCommitter(db)), jwtMiddleware, userHandler.Enable2FA)
	userGroup.Post("/disable-2fa", middleware.SetTransaction(database.NewGormCommitter(db)), jwtMiddleware, userHandler.Disable2FA)
	userGroup.Post("/logout", middleware.SetTransaction(database.NewGormCommitter(db)), jwtMiddleware, userHandler.Logout)
}

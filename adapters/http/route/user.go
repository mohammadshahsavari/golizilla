package route

import (
	"golizilla/adapters/http/handler"
	"golizilla/adapters/http/handler/middleware"
	"golizilla/config"
	"golizilla/core/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(
	app *fiber.App,
	db *gorm.DB,
	cfg *config.Config,
	userService service.IUserService,
	emailService service.IEmailService,
	roleService service.IRoleService,
) {
	// Create a group for user routes
	userGroup := app.Group("/users")

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, emailService, roleService, cfg)

	// Public routes
	userGroup.Post("/register", userHandler.CreateUser)
	userGroup.Post("/login", userHandler.Login)
	userGroup.Post("/verify-email", userHandler.VerifySignup)
	userGroup.Post("/verify-login", userHandler.VerifyLogin)

	// Initialize the JWT middleware with the config
	authMiddleware := middleware.AuthMiddleware(cfg)
	userGroup.Use(middleware.ContextMiddleware())

	// Protected routes
	userGroup.Get("/profile", authMiddleware, userHandler.GetProfile)
	userGroup.Put("/profile/update", authMiddleware, userHandler.UpdateProfile)
	userGroup.Get("/profile/notifications", authMiddleware, userHandler.GetNotificationListList)
	userGroup.Post("/enable-2fa", authMiddleware, userHandler.Enable2FA)
	userGroup.Post("/disable-2fa", authMiddleware, userHandler.Disable2FA)
	userGroup.Post("/logout", authMiddleware, userHandler.Logout)
}

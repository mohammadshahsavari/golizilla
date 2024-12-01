package route

import (
	"golizilla/config"
	"golizilla/domain/repository"
	"golizilla/handler"
	"golizilla/handler/middleware"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, database *gorm.DB, cfg *config.Config) {
	// Create a group for user routes
	userGroup := app.Group("/users")

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)

	// Initialize services
	emailService := service.NewEmailService(cfg)
	userService := service.NewUserService(userRepo, emailService)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, emailService, cfg)

	// Public routes
	userGroup.Post("/register", userHandler.CreateUser)
	userGroup.Post("/login", userHandler.Login)
	userGroup.Post("/verify-email", userHandler.VerifySignup)
	userGroup.Post("/verify-login", userHandler.VerifyLogin)

	// Initialize the JWT middleware with the config
	jwtMiddleware := middleware.AuthMiddleware(cfg)

	// Protected routes
	userGroup.Get("/profile", jwtMiddleware, userHandler.GetProfile)
	userGroup.Put("/profile/update", jwtMiddleware, userHandler.UpdateProfile)
	userGroup.Get("/profile/notifications", jwtMiddleware, userHandler.GetNotificationListList)
	userGroup.Post("/enable-2fa", jwtMiddleware, userHandler.Enable2FA)
	userGroup.Post("/disable-2fa", jwtMiddleware, userHandler.Disable2FA)
	userGroup.Post("/logout", jwtMiddleware, userHandler.Logout)
}

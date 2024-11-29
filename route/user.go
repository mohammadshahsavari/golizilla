package route

import (
	"golizilla/config"
	"golizilla/handler"
	"golizilla/handler/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userHandler *handler.UserHandler, cfg *config.Config) {
	// Create a group for user routes
	userGroup := app.Group("/users")

	// Public routes
	userGroup.Post("/register", userHandler.CreateUser)
	userGroup.Post("/login", userHandler.Login)
	userGroup.Post("/verify-email", userHandler.VerifySignup)
	userGroup.Post("/verify-login", userHandler.VerifyLogin)

	// Initialize the JWT middleware with the config
	jwtMiddleware := middleware.JWTMiddleware(cfg)

	// Protected routes
	userGroup.Get("/profile", jwtMiddleware, userHandler.GetProfile)
	userGroup.Post("/enable-2fa", jwtMiddleware, userHandler.Enable2FA)
	userGroup.Post("/disable-2fa", jwtMiddleware, userHandler.Disable2FA)
	userGroup.Post("/logout", jwtMiddleware, userHandler.Logout)
}

package route

import (
	"fmt"
	"log"
	"time"

	"golizilla/config"
	"golizilla/domain/repository"
	"golizilla/handler/middleware"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func RunServer(cfg *config.Config, database *gorm.DB) {
	// Initialize Fiber app with middleware
	app := fiber.New()

	// Add middleware for logging, panic recovery, and CORS
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Add Logging IDs to Context middleware
	middleware.InitSessionStore(cfg)
	app.Use(middleware.ContextMiddleware())

	// Add additional middleware for rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as the key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many requests. Please try again later.",
			})
		},
	}))

	// Initialize repositories
	questionRepo := repository.NewQuestionRepository(database)
	questionnariRepo := repository.NewQuestionnaireRepository(database)
	answerRepo := repository.NewAnswerRepository(database)
	roleRepo := repository.NewRoleRepository(database)
	userRepo := repository.NewUserRepository(database)
	rolePrivilegeRepo := repository.NewRolePrivilegeRepository(database)

	// Initialize services
	questionService := service.NewQuestionService(questionRepo)
	questionnariService := service.NewQuestionnaireService(questionnariRepo)
	roleService := service.NewRoleService(roleRepo, userRepo, rolePrivilegeRepo)
	authorizationsService := service.NewAuthorizationService(roleService)
	emailService := service.NewEmailService(cfg)
	userService := service.NewUserService(userRepo, emailService)
	answerService := service.NewAnswerService(answerRepo)

	// Setup routes
	SetupUserRoutes(app, database, cfg, userService, emailService, roleService)
	SetupQuestionnariRoutes(app, database, cfg, questionnariService, authorizationsService)
	SetupQuestionRoutes(app, database, cfg, questionService)
	SetupAnswerRoutes(app, database, cfg, answerService)

	// Start the server
	host := cfg.Host
	port := cfg.Port
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", host, port)))
}

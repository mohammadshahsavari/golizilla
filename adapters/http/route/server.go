package route

import (
	"fmt"
	"log"
	"time"

	"golizilla/adapters/http/handler/middleware"
	customLogger "golizilla/adapters/persistence/logger"
	"golizilla/config"
	"golizilla/core/port/repository"
	"golizilla/core/service"

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
	app.Use(customLogger.ResponseLogger(cfg))

	app.Use(recover.New())
	app.Use(cors.New())

	// Add Logging IDs to Context middleware
	middleware.InitSessionStore(cfg)
	app.Use(middleware.ContextMiddleware())
	app.Use(middleware.SetUserContext)
	app.Use(middleware.SetTransaction(database))

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
	questionnaireRepo := repository.NewQuestionnaireRepository(database)
	answerRepo := repository.NewAnswerRepository(database)
	roleRepo := repository.NewRoleRepository(database)
	userRepo := repository.NewUserRepository(database)
	rolePrivilegeRepo := repository.NewRolePrivilegeRepository(database)
	rolePrivilegeOnInstanceRepo := repository.NewRolePrivilegeOnInstanceRepository(database)
	submissionRepo := repository.NewSubmissionRepository(database)
	adminRepo := repository.NewAdminRepository(database)

	// Initialize services
	questionService := service.NewQuestionService(questionRepo)
	questionnaireService := service.NewQuestionnaireService(questionnaireRepo)
	roleService := service.NewRoleService(roleRepo, userRepo, rolePrivilegeRepo, rolePrivilegeOnInstanceRepo)
	authorizationsService := service.NewAuthorizationService(roleService)
	emailService := service.NewEmailService(cfg)
	userService := service.NewUserService(userRepo, emailService)
	answerService := service.NewAnswerService(answerRepo)
	coreService := service.NewCoreService(questionRepo, submissionRepo, questionnaireRepo, answerRepo)
	adminService := service.NewAdminService(adminRepo)

	// Setup routes
	SetupUserRoutes(app, database, cfg, userService, emailService, roleService)
	SetupQuestionnaireRoutes(app, database, cfg, questionnaireService, authorizationsService, roleService)
	SetupQuestionRoutes(app, database, cfg, questionService)
	SetupAnswerRoutes(app, database, cfg, answerService)
	SetupCoreRoutes(app, database, cfg, coreService, roleService)
	SetupAdminRoutes(app, database, cfg, adminService)

	// Start the server
	host := cfg.Host
	port := cfg.Port
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", host, port)))
}

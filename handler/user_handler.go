package handler

import (
	"strconv"

	"golizilla/domain/model"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user model.User

	// Parse request body into user model
	if err := c.BodyParser(&user); err != nil {
		return respondWithError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	// Attempt to save the user
	if err := h.UserService.CreateUser(&user); err != nil {
		// Log the error using the standard log package
		c.Context().Logger().Printf("CreateUser error: %v", err)

		// Handle specific errors such as unique constraint violation
		if isDuplicateEntryError(err) {
			return respondWithError(c, fiber.StatusConflict, "User with the given email or username already exists")
		}

		// Handle other internal errors
		return respondWithError(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	// Successfully created user
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUserByID retrieves a user by their ID
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondWithError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Attempt to fetch the user
	user, err := h.UserService.GetUserByID(uint(id))
	if err != nil {
		// Log the error using the standard log package
		c.Context().Logger().Printf("GetUserByID error: %v", err)

		if err == gorm.ErrRecordNotFound {
			return respondWithError(c, fiber.StatusNotFound, "User not found")
		}

		// Handle other internal errors
		return respondWithError(c, fiber.StatusInternalServerError, "Failed to fetch user")
	}

	// Successfully fetched user
	return c.JSON(user)
}

// Helper function to handle responses with an error status code and message
func respondWithError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"message": message,
	})
}

// Helper function to check for unique constraint violation errors
func isDuplicateEntryError(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505" // 23505 is the PostgreSQL unique violation code
	}
	return false
}

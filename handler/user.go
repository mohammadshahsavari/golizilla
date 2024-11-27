package handler

import (
	"errors"
	"golizilla/handler/presenter"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Parse request body into CreateUserRequest
	var request presenter.CreateUserRequest
	if err := c.BodyParser(&request); err != nil {
		return presenter.BadRequest(c, err)
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		return presenter.BadRequest(c, err)
	}

	// Transform the request into a domain model
	user := request.ToDomain()

	// Attempt to save the user
	if err := h.UserService.CreateUser(user); err != nil {
		// Check for PostgreSQL duplicate entry error
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == presenter.PostgresUniqueViolationCode {
			return presenter.DuplicateEntry(c, err)
		}

		// Log and respond with an internal server error
		c.Context().Logger().Printf("[CreateUser] Internal error: %v", err)
		return presenter.InternalServerError(c, err)
	}

	// Respond with the created user
	return presenter.Created(c, "User created successfully", presenter.NewUserResponse(user))
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.BadRequest(c, errors.New("invalid user ID format"))
	}

	// Attempt to fetch the user
	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presenter.NotFound(c, err)
		}

		// Log and respond with an internal server error
		c.Context().Logger().Printf("[GetUserByID] Internal error: %v", err)
		return presenter.InternalServerError(c, err)
	}

	// Respond with the fetched user
	return presenter.OK(c, "User fetched successfully", presenter.NewUserResponse(user))
}

package handler

import (
	"errors"
	"fmt"
	"golizilla/handler/presenter"
	"golizilla/service"
	"golizilla/service/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserService  service.IUserService
	EmailService service.IEmailService
}

func NewUserHandler(userService service.IUserService, emailService service.IEmailService) *UserHandler {
	return &UserHandler{
		UserService:  userService,
		EmailService: emailService,
	}
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

	// Generate email verification code and expiry
	verificationCode := utils.GenerateRandomCode(6) // We'll create this utility function
	user.EmailVerificationCode = verificationCode
	user.EmailVerificationExpiry = time.Now().Add(15 * time.Minute)
	user.IsActive = false // Ensure the user is inactive until email is verified

	// Attempt to save the user
	if err := h.UserService.CreateUser(user); err != nil {
		// Handle errors as before
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == presenter.PostgresUniqueViolationCode {
			return presenter.DuplicateEntry(c, err)
		}
		c.Context().Logger().Printf("[CreateUser] Internal error: %v", err)
		return presenter.InternalServerError(c, err)
	}

	// Send verification email
	err := h.EmailService.SendEmail(
		user.Email,
		"Email Verification",
		fmt.Sprintf("Your verification code is: %s", verificationCode),
	)
	if err != nil {
		c.Context().Logger().Printf("[CreateUser] Failed to send verification email: %v", err)
		return presenter.InternalServerError(c, errors.New("failed to send verification email"))
	}

	// Respond with a message indicating that verification is required
	return presenter.Created(c, "User created successfully. Please verify your email.", nil)
}

func (h *UserHandler) VerifySignup(c *fiber.Ctx) error {
	// Parse request body
	var request presenter.VerifyEmailRequest
	if err := c.BodyParser(&request); err != nil {
		return presenter.BadRequest(c, err)
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		return presenter.BadRequest(c, err)
	}

	// Attempt to verify the user's email
	err := h.UserService.VerifyEmail(request.Email, request.Code)
	if err != nil {
		if err == service.ErrInvalidVerificationCode {
			return presenter.BadRequest(c, errors.New("invalid or expired verification code"))
		}
		c.Context().Logger().Printf("[VerifySignup] Internal error: %v", err)
		return presenter.InternalServerError(c, err)
	}

	// Respond with success
	return presenter.OK(c, "Email verified successfully", nil)
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

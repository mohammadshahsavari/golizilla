package handler

import (
	"errors"
	"golizilla/config"
	"golizilla/domain/model"
	"golizilla/handler/presenter"
	"golizilla/internal/apperrors"
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
	Config       *config.Config
}

func NewUserHandler(userService service.IUserService, emailService service.IEmailService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		UserService:  userService,
		EmailService: emailService,
		Config:       cfg,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	ctx := c.Context()
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
	if err := h.UserService.CreateUser(ctx, user); err != nil {
		// Handle errors as before
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == presenter.PostgresUniqueViolationCode {
			return presenter.DuplicateEntry(c, err)
		}
		c.Context().Logger().Printf("[CreateUser] Internal error: %v", err)
		return presenter.InternalServerError(c, err)
	}

	// Send verification email
	emailData := map[string]interface{}{
		"Username":         user.Username,
		"VerificationCode": verificationCode,
	}

	err := h.EmailService.SendEmail(ctx,
		[]string{user.Email},
		"Email Verification",
		"verification.html",
		emailData,
	)
	if err != nil {
		c.Context().Logger().Printf("[CreateUser] Failed to send verification email: %v", err)
		return presenter.InternalServerError(c, errors.New("failed to send verification email"))
	}

	// Respond with a message indicating that verification is required
	return presenter.Created(c, "User created successfully. Please verify your email.", nil)
}

func (h *UserHandler) VerifySignup(c *fiber.Ctx) error {
	ctx := c.Context()
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
	err := h.UserService.VerifyEmail(ctx, request.Email, request.Code)
	if err != nil {
		return h.handleError(c, err)
	}

	// Respond with success
	return presenter.OK(c, "Email verified successfully", nil)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()
	// Parse request body into LoginRequest
	var request presenter.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return presenter.BadRequest(c, err)
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		return presenter.BadRequest(c, err)
	}

	// Authenticate the user
	user, err := h.UserService.AuthenticateUser(ctx, request.Email, request.Password)
	if err != nil {
		return h.handleError(c, err)
	}

	// Check if user is active
	if !user.IsActive {
		return presenter.Forbidden(c, errors.New("email not verified"))
	}

	// Check if 2FA is enabled
	if user.IsTwoFAEnabled {
		// Generate 2FA code
		twoFACode := utils.GenerateRandomCode(6)
		user.TwoFACode = twoFACode
		user.TwoFACodeExpiry = time.Now().Add(10 * time.Minute)

		// Update user in the database
		if err := h.UserService.UpdateUser(ctx, user); err != nil {
			c.Context().Logger().Printf("[Login] Failed to update user for 2FA: %v", err)
			return presenter.InternalServerError(c, err)
		}

		// Send 2FA code via email
		emailData := map[string]interface{}{
			"Username":  user.Username,
			"TwoFACode": twoFACode,
		}

		err = h.EmailService.SendEmail(ctx,
			[]string{user.Email},
			"Your 2FA Code",
			"2fa.html",
			emailData,
		)
		if err != nil {
			return h.handleError(c, err)
		}

		// Respond indicating that 2FA code has been sent
		return presenter.OK(c, "2FA code sent to your email", nil)
	}

	// If 2FA is not enabled, proceed to generate JWT token
	return h.generateAndSetToken(c, user)
}

func (h *UserHandler) VerifyLogin(c *fiber.Ctx) error {
	ctx := c.Context()
	// Parse request body into Verify2FARequest
	var request presenter.Verify2FARequest
	if err := c.BodyParser(&request); err != nil {
		return presenter.BadRequest(c, err)
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		return presenter.BadRequest(c, err)
	}

	// Find the user by email
	user, err := h.UserService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return h.handleError(c, err)
	}

	// Check if 2FA code matches and hasn't expired
	if user.TwoFACode != request.Code || time.Now().After(user.TwoFACodeExpiry) {
		err := apperrors.ErrInvalidTwoFACode
		return h.handleError(c, err)
	}

	// Clear the 2FA code fields
	user.TwoFACode = ""
	user.TwoFACodeExpiry = time.Time{}

	// Update user in the database
	if err := h.UserService.UpdateUser(ctx, user); err != nil {
		c.Context().Logger().Printf("[VerifyLogin] Failed to update user: %v", err)
		return presenter.InternalServerError(c, err)
	}

	// Generate JWT token and set cookie
	return h.generateAndSetToken(c, user)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("user_id").(uuid.UUID)

	// Fetch user details
	user, err := h.UserService.GetUserByID(ctx, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Respond with user data
	return presenter.OK(c, "User profile fetched successfully", presenter.NewUserResponse(user))
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	ctx := c.Context()
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.BadRequest(c, errors.New("invalid user ID format"))
	}

	// Attempt to fetch the user
	user, err := h.UserService.GetUserByID(ctx, id)
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

func (h *UserHandler) generateAndSetToken(c *fiber.Ctx, user *model.User) error {
	// Generate JWT token
	tokenString, err := utils.GenerateJWT(user.ID, h.Config.JWTSecretKey, h.Config.JWTExpiresIn)
	if err != nil {
		c.Context().Logger().Printf("[generateAndSetToken] Failed to generate JWT: %v", err)
		return h.handleError(c, apperrors.ErrFailedToGenerateToken)
	}

	// Set JWT token in cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(h.Config.JWTExpiresIn),
		HTTPOnly: true,
		Secure:   h.Config.Env == "production",
		SameSite: "Strict",
	})

	// Respond with success
	return presenter.OK(c, "Login successful", nil)
}

func (h *UserHandler) Enable2FA(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("user_id").(uuid.UUID)

	// Fetch user details
	user, err := h.UserService.GetUserByID(ctx, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Enable 2FA
	user.IsTwoFAEnabled = true

	// Update user in the database
	if err := h.UserService.UpdateUser(ctx, user); err != nil {
		c.Context().Logger().Printf("[Enable2FA] Failed to update user: %v", err)
		return presenter.InternalServerError(c, err)
	}

	return presenter.OK(c, "Two-factor authentication enabled", nil)
}

func (h *UserHandler) Disable2FA(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("user_id").(uuid.UUID)

	// Fetch user details
	user, err := h.UserService.GetUserByID(ctx, userID)
	if err != nil {
		c.Context().Logger().Printf("[Disable2FA] Internal error: %v", err)
		return presenter.InternalServerError(c, err)
	}

	// Disable 2FA
	user.IsTwoFAEnabled = false

	// Update user in the database
	if err := h.UserService.UpdateUser(ctx, user); err != nil {
		c.Context().Logger().Printf("[Disable2FA] Failed to update user: %v", err)
		return presenter.InternalServerError(c, err)
	}

	return presenter.OK(c, "Two-factor authentication disabled", nil)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// Clear the auth_token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiry in the past
		HTTPOnly: true,
		Secure:   h.Config.Env == "production",
		SameSite: "Strict",
	})

	return presenter.OK(c, "Logged out successfully", nil)
}

func (h *UserHandler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		return presenter.SendError(c, fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, apperrors.ErrAccountLocked):
		return presenter.SendError(c, fiber.StatusForbidden, err.Error())
	case errors.Is(err, apperrors.ErrEmailAlreadyExists):
		return presenter.SendError(c, fiber.StatusConflict, err.Error())
	case errors.Is(err, apperrors.ErrFailedToSendEmail):
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	case errors.Is(err, apperrors.ErrInvalidVerificationCode):
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	case errors.Is(err, apperrors.ErrInvalidTwoFACode):
		return presenter.SendError(c, fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, apperrors.ErrUserNotFound):
		return presenter.SendError(c, fiber.StatusNotFound, err.Error())
	case errors.Is(err, apperrors.ErrFailedToGenerateToken):
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	default:
		c.Context().Logger().Printf("Internal error: %v", err)
		return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
	}
}

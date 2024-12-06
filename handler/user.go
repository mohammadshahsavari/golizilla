package handler

import (
	"errors"
	"golizilla/config"
	"golizilla/domain/model"
	"golizilla/handler/middleware"
	"golizilla/handler/presenter"
	"golizilla/internal/apperrors"
	"golizilla/internal/logmessages"
	"golizilla/persistence/logger"
	privilegeconstants "golizilla/internal/privilegeConstants"
	"golizilla/service"
	"golizilla/service/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserHandler struct {
	UserService  service.IUserService
	EmailService service.IEmailService
	RoleService  service.IRoleService
	Config       *config.Config
}

func NewUserHandler(userService service.IUserService,
	emailService service.IEmailService,
	roleService service.IRoleService,
	cfg *config.Config) *UserHandler {
	return &UserHandler{
		UserService:  userService,
		EmailService: emailService,
		RoleService:  roleService,
		Config:       cfg,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserCreateBegin,
	})

	// Parse request body into CreateUserRequest
	var request presenter.CreateUserRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Transform the request into a domain model
	user := request.ToDomain()

	// Generate email verification code and expiry
	verificationCode := utils.GenerateRandomCode(6) // We'll create this utility function
	user.EmailVerificationCode = verificationCode
	user.EmailVerificationExpiry = time.Now().Add(h.Config.VerificationExpiresIn)
	user.IsActive = false // Ensure the user is inactive until email is verified

	// Attempt to save the user
	if err := h.UserService.CreateUser(ctx, user); err != nil {
		// Handle errors as before
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			logger.GetLogger().LogWarningFromContext(ctx, logger.LogFields{
				Service: logmessages.LogUserHandler,
				Message: apperrors.ErrEmailAlreadyExists.Error(),
			})
			return presenter.SendError(c, fiber.StatusConflict, apperrors.ErrEmailAlreadyExists.Error())
		}
		// c.Context().Logger().Printf("[CreateUser] Internal error: %v", err)
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	role, err := h.RoleService.CreateRole(ctx, user.Username, "Default")
	if err != nil {
		// c.Context().Logger().Printf("[CreateRole] Internal error: %v", err)
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	//add more privileges
	if err := h.RoleService.AddPrivilege(ctx, role.ID, privilegeconstants.CreateQuestionnari); err != nil {
		// c.Context().Logger().Printf("[CreateRole] Internal error: %v", err)
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	// Send verification email
	emailData := map[string]interface{}{
		"Username":         user.Username,
		"VerificationCode": verificationCode,
		"Expire":           h.Config.VerificationExpiresIn.Minutes(),
	}

	err = h.EmailService.SendEmail(ctx,
		[]string{user.Email},
		"Email Verification",
		"verification.html",
		emailData,
	)
	if err != nil {
		// c.Context().Logger().Printf("[CreateUser] Failed to send verification email: %v", err)
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, apperrors.ErrFailedToSendEmail)
	}

	// Respond with a message indicating that verification is required
	err = presenter.Send(c, fiber.StatusCreated, true, "User created successfully. Please verify your email.", nil, nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserCreateSuccessful,
	})

	return nil
}

func (h *UserHandler) VerifySignup(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserVerifySignupBegin,
	})

	// Parse request body
	var request presenter.VerifyEmailRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Attempt to verify the user's email
	err := h.UserService.VerifyEmail(ctx, request.Email, request.Code)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	// Respond with success
	err = presenter.Send(c, fiber.StatusOK, true, "Email verified successfully", nil, nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserVerifySignupSuccessful,
	})

	return nil
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserLoginBegin,
	})

	// Parse request body into LoginRequest
	var request presenter.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Authenticate the user
	user, err := h.UserService.AuthenticateUser(ctx, request.Email, request.Password)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	// Check if user is active
	if !user.IsActive {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: apperrors.ErrAccountLocked.Error(),
		})
		return presenter.SendError(c, fiber.StatusForbidden, apperrors.ErrAccountLocked.Error())
	}

	// Check if 2FA is enabled
	if user.IsTwoFAEnabled {
		// Generate 2FA code
		twoFACode := utils.GenerateRandomCode(6)
		user.TwoFACode = twoFACode
		user.TwoFACodeExpiry = time.Now().Add(h.Config.TwoFAExpiresIn)

		// Update user in the database
		if err := h.UserService.UpdateUser(ctx, user); err != nil {
			// c.Context().Logger().Printf("[Login] Failed to update user for 2FA: %v", err)
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogUserHandler,
				Message: err.Error(),
			})
			return h.handleError(c, err)
		}

		// Send 2FA code via email
		emailData := map[string]interface{}{
			"Username":  user.Username,
			"TwoFACode": twoFACode,
			"Expire":    h.Config.TwoFAExpiresIn.Minutes(),
		}

		err = h.EmailService.SendEmail(ctx,
			[]string{user.Email},
			"Your 2FA Code",
			"2fa.html",
			emailData,
		)
		if err != nil {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogUserHandler,
				Message: err.Error(),
			})
			return h.handleError(c, err)
		}

		// Respond indicating that 2FA code has been sent
		return presenter.Send(c, fiber.StatusOK, true, "2FA code sent to your email", nil, nil)
	}

	// Get session and set user ID
	sess, err := middleware.Store.Get(c)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, "Failed to create session")
	}

	sess.Set("user_id", user.ID.String())

	// Save session
	if err := sess.Save(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, "Failed to save session")
	}

	// If 2FA is not enabled, proceed to generate JWT token
	err = h.generateAndSetToken(c, user)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserLoginSuccessful,
	})
	
	return nil
}

func (h *UserHandler) VerifyLogin(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserVerifyLoginBegin,
	})

	// Parse request body into Verify2FARequest
	var request presenter.Verify2FARequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Find the user by email
	user, err := h.UserService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	// Check if 2FA code matches and hasn't expired
	if user.TwoFACode != request.Code || time.Now().After(user.TwoFACodeExpiry) {
		err := apperrors.ErrInvalidTwoFACode
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	// Clear the 2FA code fields
	user.TwoFACode = ""
	user.TwoFACodeExpiry = time.Time{}

	// Update user in the database
	if err := h.UserService.UpdateUser(ctx, user); err != nil {
		// c.Context().Logger().Printf("[VerifyLogin] Failed to update user: %v", err)
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	// Get session and set user ID
	sess, err := middleware.Store.Get(c)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, "Failed to create session")
	}

	sess.Set("user_id", user.ID.String())

	// Save session
	if err := sess.Save(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, "Failed to save session")
	}

	// Generate JWT token and set cookie
	err = h.generateAndSetToken(c, user)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserVerifyLoginSuccessful,
	})

	return nil
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserGetProfileBegin,
	})

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: apperrors.ErrInvalidUserID.Error(),
		})
		return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
	}

	// Fetch user details
	user, err := h.UserService.GetUserByID(ctx, userID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	// Respond with user data
	err = presenter.Send(c, fiber.StatusOK, true, "User profile fetched successfully", presenter.NewUserResponse(user), nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserGetProfileSuccessful,
	})

	return nil
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserUpdateProfileBegin,
	})

	var request presenter.UpdateProfileRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidUserDateOfBirth.Error())
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: apperrors.ErrInvalidUserID.Error(),
		})
		return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
	}

	user := request.ToDomain()
	user.ID = userID

	err := h.UserService.UpdateProfile(ctx, user)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		// c.Context().Logger().Printf("[UpdateProfile] Internal error: %v", err)
		return h.handleError(c, err)
	}

	err = presenter.Send(c, fiber.StatusOK, true, "user profile updated successfully", presenter.NewUserResponse(user), nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserUpdateProfileSuccessful,
	})

	return nil
}

func (h *UserHandler) GetNotificationListList(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserGetNotificationListBegin,
	})

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: apperrors.ErrInvalidUserID.Error(),
		})
		return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
	}

	notifList, err := h.UserService.GetNotificationList(ctx, userID)
	if err != nil {
		// c.Context().Logger().Printf("[GetNotifications] Internal error: %v", err)
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	err = presenter.Send(c, fiber.StatusOK, true, "user notifications successfully fetched", notifList, nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserGetNotificationListSuccessful,
	})

	return nil
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserGetUserByIDBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidUserID.Error())
	}

	// Attempt to fetch the user
	user, err := h.UserService.GetUserByID(ctx, id)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	// Respond with the fetched user
	err = presenter.Send(c, fiber.StatusOK, true, "User fetched successfully", presenter.NewUserResponse(user), nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserGetUserByIDSuccessful,
	})

	return nil
}

func (h *UserHandler) generateAndSetToken(c *fiber.Ctx, user *model.User) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserGenerateAndSetJWTTokenBegin,
	})

	// Generate JWT token
	tokenString, err := utils.GenerateJWT(user.ID, h.Config.JWTSecretKey, h.Config.JWTExpiresIn)
	if err != nil {
		// c.Context().Logger().Printf("[generateAndSetToken] Failed to generate JWT: %v", err)
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
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
	err = presenter.Send(c, fiber.StatusOK, true, "Login successful", nil, nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserGenerateAndSetJWTTokenSuccessful,
	})

	return nil
}

func (h *UserHandler) Enable2FA(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("user_id").(uuid.UUID)

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserEnable2FABegin,
	})

	// Fetch user details
	user, err := h.UserService.GetUserByID(ctx, userID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	// Enable 2FA
	user.IsTwoFAEnabled = true

	// Update user in the database
	if err := h.UserService.UpdateUser(ctx, user); err != nil {
		// c.Context().Logger().Printf("[Enable2FA] Failed to update user: %v", err)
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	err = presenter.Send(c, fiber.StatusOK, true, "Two-factor authentication enabled", nil, nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserEnable2FASuccessful,
	})

	return nil
}

func (h *UserHandler) Disable2FA(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("user_id").(uuid.UUID)

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserDisable2FABegin,
	})

	// Fetch user details
	user, err := h.UserService.GetUserByID(ctx, userID)
	if err != nil {
		// c.Context().Logger().Printf("[Disable2FA] Internal error: %v", err)
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	// Disable 2FA
	user.IsTwoFAEnabled = false

	// Update user in the database
	if err := h.UserService.UpdateUser(ctx, user); err != nil {
		// c.Context().Logger().Printf("[Disable2FA] Failed to update user: %v", err)
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return h.handleError(c, err)
	}

	err = presenter.Send(c, fiber.StatusOK, true, "Two-factor authentication disabled", nil, nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserDisable2FASuccessful,
	})

	return nil
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserLogoutBegin,
	})

	// Clear the auth_token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiry in the past
		HTTPOnly: true,
		Secure:   h.Config.Env == "production",
		SameSite: "Strict",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiry in the past
		HTTPOnly: true,
		Secure:   h.Config.Env == "production",
		SameSite: "Strict",
	})

	err := presenter.Send(c, fiber.StatusOK, true, "Logged out successfully", nil, nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogUserHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogUserHandler,
		Message: logmessages.LogUserLogoutBegin,
	})

	return nil
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

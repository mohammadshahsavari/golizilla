package apperrors

import "errors"

// Define error constants
var (
	ErrInvalidInput            = errors.New("invalid input")
	ErrInvalidCredentials      = errors.New("invalid email or password")
	ErrInvalidVerificationCode = errors.New("invalid or expired verification code")
	ErrAccountLocked           = errors.New("your account is locked. please try again later")
	ErrEmailAlreadyExists      = errors.New("email already in use")
	ErrNotFound                = errors.New("resource not found")
	ErrInternalServerError     = errors.New("internal server error")
	ErrInvalidTwoFACode        = errors.New("invalid or expired 2FA code")
	ErrUserNotFound            = errors.New("user not found")
	ErrFailedToGenerateToken   = errors.New("failed to generate token")
	ErrFailedToSendEmail       = errors.New("failed to send verification email")
	// Add more as needed
)

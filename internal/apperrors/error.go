package apperrors

import (
	"errors"
)

// Define error constants
var (
	ErrInvalidInput               = errors.New("invalid input")
	ErrInvalidCredentials         = errors.New("invalid email or password")
	ErrInvalidVerificationCode    = errors.New("invalid or expired verification code")
	ErrAccountLocked              = errors.New("your account is locked. please try again later")
	ErrRecordAlreadyExists        = errors.New("record already in use")
	ErrNotFound                   = errors.New("resource not found")
	ErrInternalServerError        = errors.New("internal server error")
	ErrInvalidTwoFACode           = errors.New("invalid or expired 2FA code")
	ErrUserNotFound               = errors.New("user not found")
	ErrFailedToGenerateToken      = errors.New("failed to generate token")
	ErrFailedToSendEmail          = errors.New("failed to send verification email")
	ErrMissingAuthToken           = errors.New("missing authentication token")
	ErrInvalidAuthToken           = errors.New("invalid or expired authentication token")
	ErrUnexpectedSigningMethod    = errors.New("unexpected signing method")
	ErrInvalidTokenClaims         = errors.New("invalid token claims")
	ErrInvalidUserID              = errors.New("invalid user ID in token")
	ErrInvalidUserDateOfBirth     = errors.New("invalid user date of birth")
	ErrLackOfAuthorization        = errors.New("you are not authorized to do this")
	ErrQuestionsNotFound          = errors.New("no questions available")
	ErrQuestionnaireNotFound      = errors.New("questionnaire not found")
	ErrSubmissionLimit            = errors.New("out of submission limit")
	ErrBackIsNotAllowed           = errors.New("back is not allowed")
	ErrSubmissionNotFoundQuestion = errors.New("question not found or does not match current index")
	ErrSubmissionNoQuestion       = errors.New("no current question")
	ErrSubmissionNotInProgress    = errors.New("submission not in progress")
	ErrQuestionnareExpired        = errors.New("questionnaire has expired")
	// Add more as needed
)

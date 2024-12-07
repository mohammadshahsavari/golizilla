package presenter

import (
	"errors"
	"golizilla/core/domain/model"
	"golizilla/core/service/utils"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

// CreateUserRequest defines the structure of the incoming request for creating a user.
type CreateUserRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	NationalID string `json:"national_id"`
	Password   string `json:"password"`
}

type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Verify2FARequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type UpdateProfileRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	City        string `json:"city"`
	DateOfBirth string `json:"date_of_birth"`
}

// Validate validates the CreateUserRequest fields.
func (r *CreateUserRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(r.Email) {
		return errors.New("email format is invalid") // Email validation
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.NationalID == "" {
		return errors.New("national ID is required")
	}
	if err := ValidateNationalID(r.NationalID); err != nil {
		return errors.New("national ID format is invalid")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

func (r *VerifyEmailRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Code == "" {
		return errors.New("verification code is required")
	}
	return nil
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (r *Verify2FARequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Code == "" {
		return errors.New("2FA code is required")
	}
	return nil
}

func (r *UpdateProfileRequest) Validate() error {
	if _, err := utils.ParseDate(r.DateOfBirth); err != nil {
		return errors.New("invalid date format")
	}
	return nil
}

// ToDomain transforms the CreateUserRequest into a User domain model.
func (r *CreateUserRequest) ToDomain() *model.User {
	return &model.User{
		Username:   r.Username,
		Email:      r.Email,
		Password:   r.Password,
		NationalID: r.NationalID,
	}
}

func (r *UpdateProfileRequest) ToDomain() *model.User {
	dob, _ := utils.ParseDate(r.DateOfBirth)
	return &model.User{
		FirstName:   r.FirstName,
		LastName:    r.LastName,
		City:        r.City,
		DateOfBirth: dob,
	}
}

// UserResponse defines the structure of the User object returned to the client.
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	City        string    `json:"city"`
	Wallet      uint      `json:"wallet"`
	DateOfBirth string    `json:"dateOfBirth"`
}

// NewUserResponse transforms a single User domain model into a UserResponse.
func NewUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		City:        user.City,
		Wallet:      user.Wallet,
		DateOfBirth: user.DateOfBirth.UTC().Format("2006-01-02"),
	}
}

// NewUserResponseList transforms a list of User domain models into a list of UserResponses.
func NewUserResponseList(users []*model.User) []*UserResponse {
	var response []*UserResponse
	for _, user := range users {
		response = append(response, NewUserResponse(user))
	}
	return response
}

func isValidEmail(email string) bool {
	// Simple regex for email validation
	const emailRegex = `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// isValidNationalID validates a national ID format (placeholder).
func ValidateNationalID(nationalID string) error {
	// Check length
	if len(nationalID) != 10 {
		return errors.New("national ID must be exactly 10 digits")
	}

	// Parse each digit
	sum := 0
	for i := 0; i < 9; i++ {
		digit, err := strconv.Atoi(string(nationalID[i]))
		if err != nil {
			return errors.New("national ID must contain only numeric characters")
		}
		sum += digit * (10 - i)
	}

	// Extract check digit
	checkDigit, err := strconv.Atoi(string(nationalID[9]))
	if err != nil {
		return errors.New("invalid check digit in national ID")
	}

	// Calculate modulus and validate
	sum %= 11
	if (sum < 2 && checkDigit == sum) || (sum >= 2 && checkDigit == 11-sum) {
		return nil // Valid national ID
	}

	return errors.New("invalid national ID format")
}

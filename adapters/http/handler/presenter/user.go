package presenter

import (
	"errors"
	"golizilla/core/domain/model"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	City        *string    `json:"city,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
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
		return errors.New("password must be at least 6 characters long")
	}
	if !isValidPassword(r.Password) {
		return errors.New("password must contain at least one number, one uppercase letter, and one lowercase letter")
	}
	return nil
}

func (r *VerifyEmailRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(r.Email) {
		return errors.New("email format is invalid") // Email validation
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
	if !isValidEmail(r.Email) {
		return errors.New("email format is invalid") // Email validation
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
	if !isValidEmail(r.Email) {
		return errors.New("email format is invalid") // Email validation
	}
	if r.Code == "" {
		return errors.New("2FA code is required")
	}
	return nil
}

func (r *UpdateProfileRequest) Validate() error {
	if r.DateOfBirth != nil {
		if r.DateOfBirth.After(time.Now()) {
			return errors.New("date of birth cant be in future")
		}
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
	updateFields := &model.User{}
	if r.City != nil {
		updateFields.City = *r.City
	}
	if r.FirstName != nil {
		updateFields.FirstName = *r.FirstName
	}
	if r.LastName != nil {
		updateFields.LastName = *r.LastName
	}
	if r.DateOfBirth != nil {
		updateFields.DateOfBirth = *r.DateOfBirth
	}

	return updateFields
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

// ValidateNationalID validates the Iranian national ID format.
func ValidateNationalID(nationalID string) error {
	// Trim spaces and validate length
	nationalID = strings.TrimSpace(nationalID)
	if len(nationalID) != 10 {
		return errors.New("national ID must be exactly 10 digits")
	}

	// Check for repeated digits (e.g., "0000000000", "1111111111")
	isRepeated := true
	for i := 1; i < len(nationalID); i++ {
		if nationalID[i] != nationalID[0] {
			isRepeated = false
			break
		}
	}
	if isRepeated {
		return errors.New("invalid national ID: repeated digits")
	}

	// Parse digits and calculate the checksum
	sum := 0
	for i := 0; i < 9; i++ {
		digit, err := strconv.Atoi(string(nationalID[i]))
		if err != nil {
			return errors.New("national ID must contain only numeric characters")
		}
		sum += digit * (10 - i)
	}

	// Extract and validate the check digit
	checkDigit, err := strconv.Atoi(string(nationalID[9]))
	if err != nil {
		return errors.New("invalid check digit in national ID")
	}

	// Validate using the modulus 11 rule
	mod := sum % 11
	if (mod < 2 && checkDigit == mod) || (mod >= 2 && checkDigit == 11-mod) {
		return nil // Valid national ID
	}

	return errors.New("invalid national ID format")
}

func isValidPassword(password string) bool {
	hasNumber := false
	hasUpper := false
	hasLower := false

	for _, char := range password {
		switch {
		case char >= '0' && char <= '9':
			hasNumber = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		}
	}

	return hasNumber && hasUpper && hasLower
}

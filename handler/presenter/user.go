package presenter

import (
	"errors"
	"golizilla/domain/model"
	"regexp"

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
	if !isValidNationalID(r.NationalID) {
		return errors.New("national ID format is invalid") // Placeholder for NationalID validation
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

// ToDomain transforms the CreateUserRequest into a User domain model.
func (r *CreateUserRequest) ToDomain() *model.User {
	return &model.User{
		Username:   r.Username,
		Email:      r.Email,
		Password:   r.Password,
		NationalID: r.NationalID,
	}
}

// UserResponse defines the structure of the User object returned to the client.
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

// NewUserResponse transforms a single User domain model into a UserResponse.
func NewUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
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
func isValidNationalID(nationalID string) bool {
	for i := 0; i < 10; i++ {
		if nationalID[i] < '0' || nationalID[i] > '9' {
			return false
		}
	}
	check := int(nationalID[9] - '0')
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(nationalID[i]-'0') * (10 - i)
	}
	sum %= 11
	return (sum < 2 && check == sum) || (sum >= 2 && check+sum == 11)
}

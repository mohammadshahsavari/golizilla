package presenter

import (
	"errors"
	"golizilla/domain/model"

	"github.com/google/uuid"
)

// CreateUserRequest defines the structure of the incoming request for creating a user.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailRequest struct {
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
	if r.Password == "" {
		return errors.New("password is required")
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

// ToDomain transforms the CreateUserRequest into a User domain model.
func (r *CreateUserRequest) ToDomain() *model.User {
	return &model.User{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
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

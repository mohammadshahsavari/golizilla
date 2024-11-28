package service

import (
	"errors"
	"golizilla/domain/model"
	"golizilla/domain/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserService interface {
	CreateUser(user *model.User) error
	GetUserByID(id uuid.UUID) (*model.User, error)
	VerifyEmail(email string, code string) error
	AuthenticateUser(email string, password string) (*model.User, error)
}

type UserService struct {
	UserRepo repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) CreateUser(user *model.User) error {
	return s.UserRepo.Create(user)
}

var ErrInvalidVerificationCode = errors.New("invalid or expired verification code")

func (s *UserService) VerifyEmail(email string, code string) error {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	// Check if the code matches and hasn't expired
	if user.EmailVerificationCode != code || time.Now().After(user.EmailVerificationExpiry) {
		return ErrInvalidVerificationCode
	}

	// Update user status
	user.IsActive = true
	user.EmailVerificationCode = ""
	user.EmailVerificationExpiry = time.Time{}

	return s.UserRepo.Update(user)
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *UserService) AuthenticateUser(email string, password string) (*model.User, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if password matches
	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*model.User, error) {
	return s.UserRepo.FindByID(id)
}

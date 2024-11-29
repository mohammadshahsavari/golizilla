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
	VerifyEmail(email string, code string) error
	AuthenticateUser(email string, password string) (*model.User, error)
	GetUserByID(id uuid.UUID) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user *model.User) error
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
var ErrAccountLocked = errors.New("account is locked due to multiple failed login attempts")

func (s *UserService) AuthenticateUser(email string, password string) (*model.User, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if account is locked
	if user.AccountLocked && time.Now().Before(user.AccountLockedUntil) {
		return nil, ErrAccountLocked
	}

	// Check password
	if !user.CheckPassword(password) {
		// Increment failed login attempts
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= 5 {
			user.AccountLocked = true
			user.AccountLockedUntil = time.Now().Add(3 * time.Minute)
		}
		s.UserRepo.Update(user)
		return nil, ErrInvalidCredentials
	}

	// Reset failed login attempts on successful login
	user.FailedLoginAttempts = 0
	user.AccountLocked = false
	user.AccountLockedUntil = time.Time{}
	s.UserRepo.Update(user)

	return user, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*model.User, error) {
	return s.UserRepo.FindByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.UserRepo.FindByEmail(email)
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.UserRepo.Update(user)
}

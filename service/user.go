package service

import (
	"context"
	"errors"
	"fmt"
	"golizilla/domain/model"
	"golizilla/domain/repository"
	"golizilla/internal/apperrors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*model.User, error)
	VerifyEmail(ctx context.Context, email, code string) error
	UpdateUser(ctx context.Context, user *model.User) error
}

type UserService struct {
	UserRepo     repository.IUserRepository
	EmailService IEmailService
}

func NewUserService(userRepo repository.IUserRepository, emailService IEmailService) IUserService {
	return &UserService{
		UserRepo:     userRepo,
		EmailService: emailService,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	return s.UserRepo.Create(ctx, user)
}

func (s *UserService) VerifyEmail(ctx context.Context, email string, code string) error {
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	// Check if the code matches and hasn't expired
	if user.EmailVerificationCode != code || time.Now().After(user.EmailVerificationExpiry) {
		return apperrors.ErrInvalidVerificationCode
	}

	// Update user status
	user.IsActive = true
	user.EmailVerificationCode = ""
	user.EmailVerificationExpiry = time.Time{}

	return s.UserRepo.Update(ctx, user)
}

func (s *UserService) AuthenticateUser(ctx context.Context, email string, password string) (*model.User, error) {
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	// Check if account is locked
	if user.AccountLocked && time.Now().Before(user.AccountLockedUntil) {
		return nil, apperrors.ErrAccountLocked
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// Increment failed login attempts
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= 5 {
			user.AccountLocked = true
			user.AccountLockedUntil = time.Now().Add(15 * time.Minute)
		}
		s.UserRepo.Update(ctx, user)
		return nil, apperrors.ErrInvalidCredentials
	}

	// Reset failed login attempts on successful login
	user.FailedLoginAttempts = 0
	user.AccountLocked = false
	user.AccountLockedUntil = time.Time{}
	s.UserRepo.Update(ctx, user)

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.UserRepo.FindByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.UserRepo.Update(ctx, user)
}

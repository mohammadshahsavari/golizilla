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
	// profile services
	UpdateProfile(ctx context.Context, user *model.User) error
	GetNotificationList(ctx context.Context, userId uuid.UUID) ([]*model.Notification, error)
	TransferMoney(ctx context.Context, srcEmail string, dstEmail string, amount uint) error
	CreateNotification(ctx context.Context, userId uuid.UUID, notification string) error
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

func (s *UserService) UpdateProfile(ctx context.Context, updatedUser *model.User) error {
	// Validate input
	if updatedUser == nil {
		return fmt.Errorf("updated user information must not be nil")
	}

	// Fetch the existing user
	existingUser, err := s.UserRepo.FindByID(ctx, updatedUser.ID)
	if err != nil {
		return fmt.Errorf("failed to find user with ID %s: %w", updatedUser.ID, err)
	}

	// Check if DateOfBirth can be updated
	if updatedUser.DateOfBirth != existingUser.DateOfBirth {
		// Allow changes only if the profile was created less than 24 hours ago
		if time.Since(existingUser.CreatedAt) < 24*time.Hour {
			existingUser.DateOfBirth = updatedUser.DateOfBirth
		} else {
			return fmt.Errorf("date of birth cannot be updated after 24 hours from account creation")
		}
	}

	// Update other fields
	existingUser.FirstName = updatedUser.FirstName
	existingUser.LastName = updatedUser.LastName
	existingUser.City = updatedUser.City

	// Save changes to the repository
	if err := s.UserRepo.Update(ctx, existingUser); err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	return nil
}

// it's need to test
func (s *UserService) GetNotificationList(ctx context.Context, userId uuid.UUID) ([]*model.Notification, error) {
	// Fetch the user with preloaded notifications
	user, err := s.UserRepo.FindByIDWithNotifications(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find user or notifications: %w", err)
	}

	// If the user doesn't exist
	if user == nil {
		return nil, fmt.Errorf("user with ID %s not found", userId)
	}

	return user.NotificationList, nil
}


func (s *UserService) TransferMoney(ctx context.Context, srcEmail string, dstEmail string, amount uint) error {
	// Validate input
	if srcEmail == "" || dstEmail == "" {
		return fmt.Errorf("source and destination email addresses must not be empty")
	}
	if amount == 0 {
		return fmt.Errorf("transfer amount must be greater than zero")
	}
	if srcEmail == dstEmail {
		return fmt.Errorf("source and destination email addresses cannot be the same")
	}

	// Fetch source user
	src, err := s.UserRepo.FindByEmail(ctx, srcEmail)
	if err != nil {
		return fmt.Errorf("failed to find source user: %w", err)
	}

	// Fetch destination user
	dst, err := s.UserRepo.FindByEmail(ctx, dstEmail)
	if err != nil {
		return fmt.Errorf("failed to find destination user: %w", err)
	}

	// Check if source user has sufficient balance
	if src.Wallet < amount {
		return fmt.Errorf("insufficient balance in source wallet")
	}

	// Perform transfer
	src.Wallet -= amount
	dst.Wallet += amount

	// Update source user
	if err := s.UserRepo.Update(ctx, src); err != nil {
		return fmt.Errorf("failed to update source user: %w", err)
	}

	// Update destination user
	if err := s.UserRepo.Update(ctx, dst); err != nil {
		// Rollback source user's wallet in case of failure
		src.Wallet += amount
		_ = s.UserRepo.Update(ctx, src)
		return fmt.Errorf("failed to update destination user: %w", err)
	}

	return nil
}

func (s *UserService) CreateNotification(ctx context.Context, userId uuid.UUID, notificationMsg string) error {
	return s.UserRepo.CreateNotification(ctx, userId, &model.Notification{
		Message: notificationMsg,
	})
}

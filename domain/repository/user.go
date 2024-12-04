package repository

import (
	"context"
	"errors"
	"fmt"
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	// profile
	CreateNotification(ctx context.Context, userId uuid.UUID, notification *model.Notification) error
	FindByIDWithNotifications(ctx context.Context, userId uuid.UUID) (*model.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) CreateNotification(ctx context.Context, userId uuid.UUID, notification *model.Notification) error {
	notification.UserID = userId
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

// it's need to test
func (r *UserRepository) FindByIDWithNotifications(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Preload("NotificationList").
		Where("id = ?", userId).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

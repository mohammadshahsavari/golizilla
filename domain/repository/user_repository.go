package repository

import (
	"fmt"
	"golizilla/domain/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	// Log the user data to ensure Password is set before saving
	fmt.Printf("Creating user: %+v\n", user)

	err := r.db.Create(user).Error
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
	}
	return err
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

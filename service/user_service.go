package service

import (
	"golizilla/domain/model"
	"golizilla/domain/repository"

	"github.com/google/uuid"
)

type IUserService interface {
	CreateUser(user *model.User) error
	GetUserByID(id uuid.UUID) (*model.User, error)
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

func (s *UserService) GetUserByID(id uuid.UUID) (*model.User, error) {
	return s.UserRepo.FindByID(id)
}

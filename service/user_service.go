package service

import (
	"golizilla/domain/model"
	"golizilla/domain/repository"
)

type IUserService interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
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

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.UserRepo.FindByID(id)
}

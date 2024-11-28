package service

import (
	"golizilla/domain/model"
	"golizilla/domain/repository"
)

type RoleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) CreateRole(name, description string) (*model.Role, error) {
	role := &model.Role{
		Name:        name,
		Description: description,
	}
	if err := s.roleRepo.CreateRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) GetAllRoles() ([]*model.Role, error) {
	return s.roleRepo.GetAllRoles()
}

package service

import (
	"golizilla/domain/model"
	"golizilla/domain/repository"
	"time"
)

type RBACService interface {
	CreateRole(name, description string) (*model.Role, error)
	AssignPermissionToRole(roleID, permissionID uint) error
	GetAllRoles() ([]model.Role, error)
	AssignRoleToUser(userID, roleID uint, expiry *time.Time) error
	GetUserRoles(userID uint) ([]model.Role, error)
}

type rbacService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	userRoleRepo   repository.UserRoleRepository
}

func NewRBACService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	userRoleRepo repository.UserRoleRepository,
) RBACService {
	return &rbacService{roleRepo, permissionRepo, userRoleRepo}
}

func (s *rbacService) CreateRole(name, description string) (*model.Role, error) {
	role := &model.Role{Name: name, Description: description}
	err := s.roleRepo.Create(role)
	return role, err
}

func (s *rbacService) AssignPermissionToRole(roleID, permissionID uint) error {
	return s.roleRepo.AssignPermission(roleID, permissionID)
}

func (s *rbacService) GetAllRoles() ([]model.Role, error) {
	return s.roleRepo.List()
}

func (s *rbacService) AssignRoleToUser(userID, roleID uint, expiry *time.Time) error {
	return s.userRoleRepo.AssignRole(userID, roleID, expiry)
}

func (s *rbacService) GetUserRoles(userID uint) ([]model.Role, error) {
	return s.userRoleRepo.GetUserRoles(userID)
}

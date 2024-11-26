package repository

import (
	"golizilla/domain/model"
	"time"
)

type RoleRepository interface {
	Create(role *model.Role) error
	GetByID(id uint) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	List() ([]model.Role, error)
	AssignPermission(roleID uint, permissionID uint) error
}

type PermissionRepository interface {
	Create(permission *model.Permission) error
	GetByID(id uint) (*model.Permission, error)
	GetByName(name string) (*model.Permission, error)
	List() ([]model.Permission, error)
}

type UserRoleRepository interface {
	AssignRole(userID uint, roleID uint, expiry *time.Time) error
	GetUserRoles(userID uint) ([]model.Role, error)
	RemoveRole(userID uint, roleID uint) error
}

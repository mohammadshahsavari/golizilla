package repository

import "golizilla/domain/model"

// RoleRepository defines methods for managing roles.
type RoleRepository interface {
	CreateRole(role *model.Role) error
	GetRoleByID(id uint) (*model.Role, error)
	GetAllRoles() ([]*model.Role, error)
	DeleteRole(id uint) error
}

// PermissionRepository defines methods for managing permissions.
type PermissionRepository interface {
	CreatePermission(permission *model.Permission) error
	GetPermissionByID(id uint) (*model.Permission, error)
	GetAllPermissions() ([]*model.Permission, error)
	DeletePermission(id uint) error
}

// RolePermissionRepository defines methods for managing role-permission mappings.
type RolePermissionRepository interface {
	AssignPermissionToRole(roleID uint, permissionID uint) error
	RevokePermissionFromRole(roleID uint, permissionID uint) error
	GetPermissionsByRole(roleID uint) ([]*model.Permission, error)
}

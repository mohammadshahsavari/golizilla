package model

import "time"

// Role represents a user role with specific permissions.
type Role struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission represents a specific operation a role can perform.
type Permission struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RolePermission associates roles with their permissions.
type RolePermission struct {
	RoleID       uint `json:"role_id"`
	PermissionID uint `json:"permission_id"`
}

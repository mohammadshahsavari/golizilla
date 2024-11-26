package model

import "time"

type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"unique;not null" json:"name"` // e.g., "Super Admin", "Owner"
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions"`
}

type Permission struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"unique;not null" json:"name"` // e.g., "CREATE_QUESTIONNAIRE"
	Description string `json:"description"`
}

type UserRole struct {
	UserID uint       `gorm:"primaryKey" json:"user_id"`
	RoleID uint       `gorm:"primaryKey" json:"role_id"`
	Expiry *time.Time `json:"expiry"` // Optional: Time-limited roles
}

package model

import "github.com/google/uuid"

type RolePrivilege struct {
	RoleId      uuid.UUID `gorm:"primaryKey"`
	PrivilegeId string    `gorm:"primaryKey"`
	Role        Role      `gorm:"foreinKey:RoleId"`
	Privilege   Privilege `gorm:"foreinKey:PrivilegeId"`
}

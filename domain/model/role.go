package model

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name           string
	Description    string
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time
	Users          []*User          `gorm:"foreinKey:RoleId"`
	Privileges     []*Privilege     `gorm:"many2many:RolePrivilege"`
	Questionnaires []*Questionnaire `gorm:"many2many:RolePrivilegeOnQuestionnaire"`
}

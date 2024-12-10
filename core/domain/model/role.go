package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name           string
	Description    string
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time
	Users          []*User          `gorm:"foreinKey:RoleId"`
	Privileges     []*Privilege     `gorm:"many2many:RolePrivilege"`
	Questionnaires []*Questionnaire `gorm:"many2many:RolePrivilegeOnInstance"`
}

// BeforeCreate is a GORM hook to generate a UUID before creating a new record.
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

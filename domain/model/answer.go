package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Answer struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	QuestionID   uuid.UUID `gorm:"not null"`
	Answer       string
	AnswerOption uint
	//answer time can be added
}

// BeforeCreate is a GORM hook to generate a UUID before creating a new record.
func (a *Answer) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

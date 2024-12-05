package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;"`
	QuestionnaireId uuid.UUID `gorm:"not null"`
	Index           uint
	QuestionText    string
	Descriptive     bool
	OptionsCount    uint
	CorrectOption   uint
	MetaDataPath    string
	OptionsText     string 
	SelectedOption  uint
	Answers         []*Answer `gorm:"foreignKey:QuestionID"`
}

// BeforeCreate is a GORM hook to generate a UUID before creating a new record.
func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}


package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	ID              uuid.UUID     `gorm:"type:uuid;primary_key;"`
	QuestionnaireId uuid.UUID     `gorm:"type:uuid;not null"` // FK to Questionnaire
	Questionnaire   Questionnaire `gorm:"foreignKey:QuestionnaireId"`

	Index        uint
	QuestionText string
	Descriptive  bool
	MetaDataPath string

	// For correct option, store an ID
	CorrectOptionID *uuid.UUID
	// CorrectOption   *Option `gorm:"foreignKey:CorrectOptionID"`

	// Multiple options
	Options []Option `gorm:"foreignKey:QuestionID"`

	// Multiple answers for this question
	Answers []Answer `gorm:"foreignKey:QuestionID"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

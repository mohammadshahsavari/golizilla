package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Answer struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;"`
	QuestionID       uuid.UUID      `gorm:"type:uuid;not null"` // Foreign key to Question
	Question         Question       `gorm:"foreignKey:QuestionID"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null"` // Foreign key to User
	User             User           `gorm:"foreignKey:UserID"`
	UserSubmissionID uuid.UUID      `gorm:"type:uuid;not null"` // Foreign key to UserSubmission
	UserSubmission   UserSubmission `gorm:"foreignKey:UserSubmissionID"`

	Descriptive bool
	Text        *string
	OptionID    *uuid.UUID // Optional if answer references a chosen option
	Option      *Option    `gorm:"foreignKey:OptionID"`
}

func (a *Answer) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

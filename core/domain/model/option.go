package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Option struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null"` // FK back to Question
	Index      uint
	Text       string
}

func (o *Option) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

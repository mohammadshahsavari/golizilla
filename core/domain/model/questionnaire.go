package model

import (
	"time"

	"github.com/google/uuid"
)

type Questionnaire struct {
	Id                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerId            uuid.UUID `gorm:"not null"`
	CreatedTime        time.Time `gorm:"not null"`
	StartTime          time.Time `gorm:"not null"`
	EndTime            time.Time `gorm:"not null"`
	Random             bool
	BackCompatible     bool
	Title              string
	AnswerTime         time.Duration `gorm:"not null"`
	ParticipationCount uint
	Anonymous          bool
	SubmitLimit        uint
	Owner              User `gorm:"foreinKey:OwnerId"`
}

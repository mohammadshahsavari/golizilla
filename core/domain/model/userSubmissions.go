package model

import (
	"time"

	"github.com/google/uuid"
)

type SubmissionStatus string

const (
	SubmissionsStatusInProgress SubmissionStatus = "in_progress"
	SubmissionsStatusDone       SubmissionStatus = "done"
	SubmissionsStatusCancelled  SubmissionStatus = "cancelled"
)

type UserSubmission struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key;"`
	UserId          uuid.UUID        `gorm:"type:uuid;not null"` // FK to User
	User            User             `gorm:"foreignKey:UserId"`
	QuestionnaireId uuid.UUID        `gorm:"type:uuid;not null"` // FK to Questionnaire
	Questionnaire   Questionnaire    `gorm:"foreignKey:QuestionnaireId"`
	Status          SubmissionStatus `gorm:"not null;default:'in_progress'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// One submission has multiple answers
	// We linked it above in Answer with UserSubmissionID
	Answers []Answer `gorm:"foreignKey:UserSubmissionID"`

	// If you need to track current question index
	CurrentQuestionIndex int `gorm:"default:0"`
}

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

// User represents a user in the domain layer.
type UserSubmission struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key;"`
	UserId          uuid.UUID        `gorm:"not null"`
	QuestionnaireId uuid.UUID        `gorm:"not null"`
	AnswersId       uuid.UUID        `gorm:"not null"`
	Status          SubmissionStatus `gorm:"not null;type:submission_status;default:'in_progress'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	User            User          `gorm:"foreinKey:UserId"`
	Questionnaire   Questionnaire `gorm:"foreinKey:QuestionnaireId"`
	Answers         []Answer      `gorm:"foreignKey:AnswersID"`
}

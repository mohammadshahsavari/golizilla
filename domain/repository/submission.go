package repository

import (
	"golizilla/domain/model"

	"github.com/google/uuid"
)

type ISubmissionRepository interface {
	GetSubmissionByID(submission uuid.UUID) (*model.UserSubmission, error)
	CreateSubmission(submission *model.UserSubmission) error
	UpdateSubmission(submission *model.UserSubmission) error

	GetActiveSubmissionByUserID(userID uuid.UUID) (*model.UserSubmission, error)
}

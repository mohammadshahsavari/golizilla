package repository

import (
	"context"
	"golizilla/domain/model"

	"github.com/google/uuid"
)

type ISubmissionRepository interface {
	GetSubmissionByID(ctx context.Context, userCtx context.Context, submission uuid.UUID) (*model.UserSubmission, error)
	CreateSubmission(ctx context.Context, userCtx context.Context, submission *model.UserSubmission) error
	UpdateSubmission(ctx context.Context, userCtx context.Context, submission *model.UserSubmission) error

	GetActiveSubmissionByUserID(ctx context.Context, userCtx context.Context, userID uuid.UUID) (*model.UserSubmission, error)
}

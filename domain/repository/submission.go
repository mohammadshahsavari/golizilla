package repository

import (
	"context"
	"golizilla/domain/model"
	appContext "golizilla/handler/context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ISubmissionRepository interface {
	GetSubmissionByID(ctx context.Context, userCtx context.Context, submissionID uuid.UUID) (*model.UserSubmission, error)
	CreateSubmission(ctx context.Context, userCtx context.Context, submission *model.UserSubmission) error
	UpdateSubmission(ctx context.Context, userCtx context.Context, submission *model.UserSubmission) error
	GetActiveSubmissionByUserIDAndQuestionnaire(ctx context.Context, userCtx context.Context, userID, questionnaireID uuid.UUID) (*model.UserSubmission, error)
	// Add any other needed methods, e.g., to get the current question index, etc.
}

type SubmissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) ISubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) GetSubmissionByID(ctx context.Context, userCtx context.Context, submissionID uuid.UUID) (*model.UserSubmission, error) {
	db := appContext.GetDB(userCtx)
	if db == nil {
		db = r.db
	}
	var sub model.UserSubmission
	if err := db.WithContext(ctx).Where("id = ?", submissionID).Preload("Answers").First(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubmissionRepository) CreateSubmission(ctx context.Context, userCtx context.Context, submission *model.UserSubmission) error {
	db := appContext.GetDB(userCtx)
	if db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Create(submission).Error
}

func (r *SubmissionRepository) UpdateSubmission(ctx context.Context, userCtx context.Context, submission *model.UserSubmission) error {
	db := appContext.GetDB(userCtx)
	if db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Save(submission).Error
}

func (r *SubmissionRepository) GetActiveSubmissionByUserIDAndQuestionnaire(ctx context.Context, userCtx context.Context, userID, questionnaireID uuid.UUID) (*model.UserSubmission, error) {
	db := appContext.GetDB(userCtx)
	if db == nil {
		db = r.db
	}
	var sub model.UserSubmission
	err := db.WithContext(ctx).
		Where("user_id = ? AND questionnaire_id = ? AND status = ?", userID, questionnaireID, model.SubmissionsStatusInProgress).
		Preload("Answers").First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

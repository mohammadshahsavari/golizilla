package repository

import (
	"context"
	"errors"
	"fmt"
	myContext "golizilla/adapters/http/handler/context"
	"golizilla/core/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAnswerRepository interface {
	Create(ctx context.Context, userCtx context.Context, answer *model.Answer) (uuid.UUID, error)
	Update(ctx context.Context, userCtx context.Context, answer *model.Answer) error
	Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Answer, error)
}

type AnswerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) IAnswerRepository {
	return &AnswerRepository{
		db: db,
	}
}

func (r *AnswerRepository) Create(ctx context.Context, userCtx context.Context, answer *model.Answer) (uuid.UUID, error) {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}

	// Check if a record with the same QuestionID and SubmissionID exists
	var existingAnswer model.Answer
	err := db.WithContext(ctx).
		Where("question_id = ? AND user_submission_id = ?", answer.QuestionID, answer.UserSubmissionID).
		First(&existingAnswer).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, fmt.Errorf("failed to query existing answer: %w", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No existing record, create a new one
		result := db.WithContext(ctx).Create(answer)
		if result.Error != nil {
			return uuid.Nil, fmt.Errorf("failed to create answer: %w", result.Error)
		}
	} else {
		// Update the existing record
		existingAnswer.Descriptive = answer.Descriptive
		existingAnswer.Text = answer.Text // Update any other fields as needed
		result := db.WithContext(ctx).Save(&existingAnswer)
		if result.Error != nil {
			return uuid.Nil, fmt.Errorf("failed to update answer: %w", result.Error)
		}
		return existingAnswer.ID, nil
	}

	return answer.ID, nil
}

func (r *AnswerRepository) Update(ctx context.Context, userCtx context.Context, answer *model.Answer) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Save(answer).Error
}

func (r *AnswerRepository) Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Delete(&model.Answer{}, id).Error
}

func (r *AnswerRepository) GetByID(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Answer, error) {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	var answer model.Answer
	err := db.WithContext(ctx).Where("id = ?", id).First(&answer).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find answer by ID: %v, %w", id, err)
	}
	return &answer, nil
}

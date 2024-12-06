package repository

import (
	"context"
	"fmt"
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IQuestionRepository interface {
	Create(ctx context.Context, question *model.Question) (uuid.UUID, error)
	Update(ctx context.Context, question *model.Question) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Question, error)
}

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) IQuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (r *QuestionRepository) Create(ctx context.Context, question *model.Question) (uuid.UUID, error) {
	result := r.db.WithContext(ctx).Create(question)
	if result.Error != nil {
		return uuid.Nil, fmt.Errorf("failed to create question: %w", result.Error)
	}
	return question.ID, nil
}

func (r *QuestionRepository) Update(ctx context.Context, question *model.Question) error {
	return r.db.WithContext(ctx).Where("id = ?", question.ID).Updates(question).Error
}

func (r *QuestionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Question{}, id).Error
}

func (r *QuestionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Question, error) {
	var question model.Question
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&question).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find question by ID: %v, %w", id, err)
	}
	return &question, nil
}

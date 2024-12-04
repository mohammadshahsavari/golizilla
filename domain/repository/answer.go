package repository

import (
	"context"
	"fmt"
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAnswerRepository interface {
	Create(ctx context.Context, answer *model.Answer) (uuid.UUID, error)
	Update(ctx context.Context, answer *model.Answer) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Answer, error)
}

type AnswerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) IAnswerRepository {
	return &AnswerRepository{
		db: db,
	}
}

func (r *AnswerRepository) Create(ctx context.Context, answer *model.Answer) (uuid.UUID, error) {
	result := r.db.WithContext(ctx).Create(answer)
	if result.Error != nil {
		return uuid.Nil, fmt.Errorf("failed to create answer: %w", result.Error)
	}
	return answer.ID, nil
}

func (r *AnswerRepository) Update(ctx context.Context, answer *model.Answer) error {
	return r.db.WithContext(ctx).Save(answer).Error
}

func (r *AnswerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Answer{}, id).Error
}

func (r *AnswerRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Answer, error) {
	var answer model.Answer
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&answer).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find answer by ID: %v, %w", id, err)
	}
	return &answer, nil
}

package repository

import (
	"context"
	"fmt"
	"golizilla/domain/model"
	myContext "golizilla/handler/context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAnswerRepository interface {
	Create(ctx context.Context, userCtx context.Context, answer *model.Answer) (uuid.UUID, error)
	Update(ctx context.Context, userCtx context.Context, answer *model.Answer) error
	Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Answer, error)
	GetByQuestionId(ctx context.Context, userCtx context.Context, questionId uuid.UUID) ([]model.Answer, error)
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
	result := db.WithContext(ctx).Create(answer)
	if result.Error != nil {
		return uuid.Nil, fmt.Errorf("failed to create answer: %w", result.Error)
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

func (r *AnswerRepository) GetByQuestionId(ctx context.Context, userCtx context.Context, questionId uuid.UUID) ([]model.Answer, error) {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	var answers []model.Answer
	err := db.WithContext(ctx).Where("question_id = ?", questionId).Find(&answers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find answers by ID: %v, %w", questionId, err)
	}
	return answers, nil
}

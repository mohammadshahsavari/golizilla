package repository

import (
	"context"
	"errors"
	"fmt"
	myContext "golizilla/adapters/http/handler/context"
	"golizilla/core/domain/model"
	"golizilla/internal/apperrors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IQuestionRepository interface {
	Create(ctx context.Context, userCtx context.Context, question *model.Question) (uuid.UUID, error)
	Update(ctx context.Context, userCtx context.Context, question *model.Question) error
	Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Question, error)
	GetByQuestionnaireID(ctx context.Context, userCtx context.Context, questionnaireID uuid.UUID) ([]*model.Question, error)
	GetFullByQuestionnaireID(ctx context.Context, userCtx context.Context, questionnaireID uuid.UUID) ([]*model.Question, error)
}

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) IQuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (r *QuestionRepository) Create(ctx context.Context, userCtx context.Context, question *model.Question) (uuid.UUID, error) {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	result := db.WithContext(ctx).Create(question)
	if result.Error != nil {
		return uuid.Nil, fmt.Errorf("failed to create question: %w", result.Error)
	}
	return question.ID, nil
}

func (r *QuestionRepository) Update(ctx context.Context, userCtx context.Context, question *model.Question) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Where("id = ?", question.ID).Updates(question).Error
}

func (r *QuestionRepository) Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Delete(&model.Question{}, id).Error
}

func (r *QuestionRepository) GetByID(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Question, error) {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	var question model.Question
	if err := r.db.WithContext(userCtx).Preload("Options").First(&question, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return &question, nil
}

func (r *QuestionRepository) GetByQuestionnaireID(ctx context.Context, userCtx context.Context, questionnaireID uuid.UUID) ([]*model.Question, error) {
	db := myContext.GetDB(userCtx)
	if db == nil {
		db = r.db
	}

	var questions []*model.Question
	if err := db.WithContext(ctx).Where("questionnaire_id = ?", questionnaireID).Order("index ASC").Find(&questions).Error; err != nil {
		return nil, err
	}

	return questions, nil
}

func (r *QuestionRepository) GetFullByQuestionnaireID(ctx context.Context, userCtx context.Context, questionnaireID uuid.UUID) ([]*model.Question, error) {
	db := myContext.GetDB(userCtx)
	if db == nil {
		db = r.db
	}

	var questions []*model.Question
	if err := db.WithContext(ctx).Preload("Answers").Where("questionnaire_id = ?", questionnaireID).Order("index ASC").Find(&questions).Error; err != nil {
		return nil, err
	}

	return questions, nil
}

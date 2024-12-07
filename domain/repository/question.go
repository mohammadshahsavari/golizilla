package repository

import (
	"context"
	"fmt"
	"golizilla/domain/model"
	myContext "golizilla/handler/context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IQuestionRepository interface {
	Create(ctx context.Context, userCtx context.Context, question *model.Question) (uuid.UUID, error)
	Update(ctx context.Context, userCtx context.Context, question *model.Question) error
	Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Question, error)
	GetByQuestionnariId(ctx context.Context, userCtx context.Context, questionnariId uuid.UUID) ([]model.Question, error)
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
	err := db.WithContext(ctx).Where("id = ?", id).First(&question).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find question by ID: %v, %w", id, err)
	}
	return &question, nil
}

func (r *QuestionRepository) GetByQuestionnariId(ctx context.Context, userCtx context.Context, questionnariId uuid.UUID) ([]model.Question, error) {
	var db *gorm.DB
	if db = myContext.GetDB(userCtx); db == nil {
		db = r.db
	}
	var questions []model.Question
	err := db.WithContext(ctx).Where("questionnaire_id = ?", questionnariId).Find(&questions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find questions by questionnari Id: %v, %w", questionnariId, err)
	}
	return questions, nil
}

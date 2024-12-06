package repository

import (
	"context"
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IQuestionnaireRepository interface {
	Add(ctx context.Context, questionnaire *model.Questionnaire) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, questionnaire map[string]interface{}) error
	GetById(ctx context.Context, id uuid.UUID) (*model.Questionnaire, error)
	GetByOwnerId(ctx context.Context, ownerId uuid.UUID) ([]model.Questionnaire, error)
}

type questionnaireRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRepository(db *gorm.DB) IQuestionnaireRepository {
	return &questionnaireRepository{
		db: db,
	}
}

func (r *questionnaireRepository) Add(ctx context.Context, questionnaire *model.Questionnaire) error {
	err := r.db.WithContext(ctx).Create(questionnaire).Error
	if err != nil {
		//log
	}
	return err
}

func (r *questionnaireRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&model.Questionnaire{}, id).Error
	if err != nil {
		//log
	}
	return err
}

func (r *questionnaireRepository) Update(ctx context.Context, id uuid.UUID, questionnaire map[string]interface{}) error {
	err := r.db.WithContext(ctx).Model(&model.Questionnaire{}).Where("id = ?", id).Updates(questionnaire).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *questionnaireRepository) GetById(ctx context.Context, id uuid.UUID) (*model.Questionnaire, error) {
	var questionnaire model.Questionnaire
	err := r.db.WithContext(ctx).Find(&questionnaire, id).Error
	if err != nil {
		//log
	}

	return &questionnaire, err
}

func (r *questionnaireRepository) GetByOwnerId(ctx context.Context, ownerId uuid.UUID) ([]model.Questionnaire, error) {
	var questionnaires []model.Questionnaire
	err := r.db.WithContext(ctx).Where("owner_id = ?", ownerId).Find(&questionnaires).Error
	if err != nil {
		//log
	}

	return questionnaires, err
}

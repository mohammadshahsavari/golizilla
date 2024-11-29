package repository

import (
	"golizilla/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IQuestionnaireRepository interface {
	Add(questionnaire *model.Questionnaire) error
	Delete(id uuid.UUID) error
	Update(questionnaire *model.Questionnaire) error
	GetById(id uuid.UUID) (*model.Questionnaire, error)
	GetByOwnerId(ownerId uuid.UUID) ([]model.Questionnaire, error)
}

type questionnaireRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRepository(db *gorm.DB) IQuestionnaireRepository {
	return &questionnaireRepository{
		db: db,
	}
}

func (r *questionnaireRepository) Add(questionnaire *model.Questionnaire) error {

	err := r.db.Create(questionnaire).Error
	if err != nil {
		//log
	}
	return err
}

func (r *questionnaireRepository) Delete(id uuid.UUID) error {
	err := r.db.Delete(&model.Questionnaire{}, id).Error
	if err != nil {
		//log
	}
	return err
}

func (r *questionnaireRepository) Update(questionnaire *model.Questionnaire) error {
	err := r.db.Save(questionnaire).Error
	if err != nil {
		//log
	}
	return err
}

func (r *questionnaireRepository) GetById(id uuid.UUID) (*model.Questionnaire, error) {
	var questionnaire model.Questionnaire
	err := r.db.Find(&questionnaire, id).Error
	if err != nil {
		//log
	}

	return &questionnaire, err
}

func (r *questionnaireRepository) GetByOwnerId(ownerId uuid.UUID) ([]model.Questionnaire, error) {
	var questionnaires []model.Questionnaire
	err := r.db.Where("owner_id = ?", ownerId).Find(&questionnaires, ownerId).Error
	if err != nil {
		//log
	}

	return questionnaires, err
}

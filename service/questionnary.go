package service

import (
	"golizilla/domain/model"
	respository "golizilla/domain/repository"

	"github.com/google/uuid"
)

type IQuestionnaireService interface {
	Create(questionary *model.Questionnaire) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	Update(questionnaire *model.Questionnaire) error
	GetById(id uuid.UUID) (*model.Questionnaire, error)
	GetByOwnerId(ownerId uuid.UUID) ([]model.Questionnaire, error)
}

type questionnaireService struct {
	repo respository.IQuestionnaireRepository
}

func NewQuestionnaireService(repo repository.IQuestionnaireRepository) IQuestionnaireService {
	return &questionnaireService{
		repo: repo,
	}
}

func (q *questionnaireService) Create(questionary *model.Questionnaire) (uuid.UUID, error) {
	questionary.Id = uuid.New()
	err := q.repo.Add(questionary)
	if err != nil {
		//log
		questionary.Id = uuid.Nil
	}

	return questionary.Id, err
}

func (q *questionnaireService) Delete(id uuid.UUID) error {
	return q.repo.Delete(id)
}

func (q *questionnaireService) Update(questionnaire *model.Questionnaire) error {
	return q.repo.Update(questionnaire)
}

func (q *questionnaireService) GetById(id uuid.UUID) (*model.Questionnaire, error) {
	return q.repo.GetById(id)
}

func (q *questionnaireService) GetByOwnerId(ownerId uuid.UUID) ([]model.Questionnaire, error) {
	return q.repo.GetByOwnerId(ownerId)
}

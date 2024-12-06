package service

import (
	"context"
	"golizilla/domain/model"
	respository "golizilla/domain/repository"

	"github.com/google/uuid"
)

type IQuestionnaireService interface {
	Create(ctx context.Context, questionnaire *model.Questionnaire) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, questionnaire map[string]interface{}) error
	GetById(ctx context.Context, id uuid.UUID) (*model.Questionnaire, error)
	GetByOwnerId(ctx context.Context, ownerId uuid.UUID) ([]model.Questionnaire, error)
}

type questionnaireService struct {
	repo respository.IQuestionnaireRepository
}

func NewQuestionnaireService(repo respository.IQuestionnaireRepository) IQuestionnaireService {
	return &questionnaireService{
		repo: repo,
	}
}

func (q *questionnaireService) Create(ctx context.Context, questionnaire *model.Questionnaire) (uuid.UUID, error) {
	questionnaire.Id = uuid.New()
	err := q.repo.Add(ctx, questionnaire)
	if err != nil {
		//log
		questionnaire.Id = uuid.Nil
	}

	return questionnaire.Id, err
}

func (q *questionnaireService) Delete(ctx context.Context, id uuid.UUID) error {
	return q.repo.Delete(ctx, id)
}

func (q *questionnaireService) Update(ctx context.Context, id uuid.UUID, updateFields map[string]interface{}) error {
	return q.repo.Update(ctx, id, updateFields)
}

func (q *questionnaireService) GetById(ctx context.Context, id uuid.UUID) (*model.Questionnaire, error) {
	return q.repo.GetById(ctx, id)
}

func (q *questionnaireService) GetByOwnerId(ctx context.Context, ownerId uuid.UUID) ([]model.Questionnaire, error) {
	return q.repo.GetByOwnerId(ctx, ownerId)
}

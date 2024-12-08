package service

import (
	"context"
	"golizilla/core/domain/model"
	respository "golizilla/core/port/repository"

	"github.com/google/uuid"
)

type IQuestionnaireService interface {
	Create(ctx context.Context, userCtx context.Context, questionnaire *model.Questionnaire) (uuid.UUID, error)
	Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error
	Update(ctx context.Context, userCtx context.Context, id uuid.UUID, questionnaire map[string]interface{}) error
	GetById(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Questionnaire, error)
	GetByOwnerId(ctx context.Context, userCtx context.Context, ownerId uuid.UUID) ([]model.Questionnaire, error)
	IsOwner(ctx context.Context, userCtx context.Context, userId uuid.UUID, questionnariId uuid.UUID) (bool, error)
}

type questionnaireService struct {
	repo respository.IQuestionnaireRepository
}

func NewQuestionnaireService(repo respository.IQuestionnaireRepository) IQuestionnaireService {
	return &questionnaireService{
		repo: repo,
	}
}

func (q *questionnaireService) Create(ctx context.Context, userCtx context.Context, questionnaire *model.Questionnaire) (uuid.UUID, error) {
	questionnaire.Id = uuid.New()
	err := q.repo.Add(ctx, userCtx, questionnaire)
	if err != nil {
		//log
		questionnaire.Id = uuid.Nil
	}

	return questionnaire.Id, err
}

func (q *questionnaireService) Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error {
	return q.repo.Delete(ctx, userCtx, id)
}

func (q *questionnaireService) Update(ctx context.Context, userCtx context.Context, id uuid.UUID, updateFields map[string]interface{}) error {
	return q.repo.Update(ctx, userCtx, id, updateFields)
}

func (q *questionnaireService) GetById(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Questionnaire, error) {
	return q.repo.GetById(ctx, userCtx, id)
}

func (q *questionnaireService) GetByOwnerId(ctx context.Context, userCtx context.Context, ownerId uuid.UUID) ([]model.Questionnaire, error) {
	return q.repo.GetByOwnerId(ctx, userCtx, ownerId)
}

func (q *questionnaireService) IsOwner(ctx context.Context, userCtx context.Context, userId uuid.UUID, questionnariId uuid.UUID) (bool, error) {
	return q.repo.IsOwner(ctx, userCtx, userId, questionnariId)
}

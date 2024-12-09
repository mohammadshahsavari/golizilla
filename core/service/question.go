package service

import (
	"context"
	"golizilla/core/domain/model"
	"golizilla/core/port/repository"

	"github.com/google/uuid"
)

type IQuestionService interface {
	Create(ctx context.Context, userCtx context.Context, question *model.Question) (uuid.UUID, error)
	Update(ctx context.Context, userCtx context.Context, question *model.Question) error
	Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Question, error)
	GetFullByQuestionnaireID(ctx context.Context, userCtx context.Context, questionnaireID uuid.UUID) ([]*model.Question, error)
}

type QuestionService struct {
	QuestionRepo repository.IQuestionRepository
}

func NewQuestionService(repo repository.IQuestionRepository) IQuestionService {
	return &QuestionService{
		QuestionRepo: repo,
	}
}

func (s *QuestionService) Create(ctx context.Context, userCtx context.Context, question *model.Question) (uuid.UUID, error) {
	return s.QuestionRepo.Create(ctx, userCtx, question)
}

func (s *QuestionService) Update(ctx context.Context, userCtx context.Context, question *model.Question) error {
	return s.QuestionRepo.Update(ctx, userCtx, question)
}

func (s *QuestionService) Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error {
	return s.QuestionRepo.Delete(ctx, userCtx, id)
}

func (s *QuestionService) GetByID(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Question, error) {
	return s.QuestionRepo.GetByID(ctx, userCtx, id)
}

func (s *QuestionService) GetFullByQuestionnaireID(ctx context.Context, userCtx context.Context, questionnaireID uuid.UUID) ([]*model.Question, error) {
	return s.QuestionRepo.GetFullByQuestionnaireID(ctx, userCtx, questionnaireID)
}

package service

import (
	"context"
	"golizilla/domain/model"
	"golizilla/domain/repository"

	"github.com/google/uuid"
)

type IQuestionService interface {
	Create(ctx context.Context, question *model.Question) (uuid.UUID, error)
	Update(ctx context.Context, question *model.Question) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Question, error)
}

type QuestionService struct {
	QuestionRepo repository.IQuestionRepository
}

func NewQuestionService(repo repository.IQuestionRepository) IQuestionService {
	return &QuestionService{
		QuestionRepo: repo,
	}
}

func (s *QuestionService) Create(ctx context.Context, question *model.Question) (uuid.UUID, error) {
	return s.QuestionRepo.Create(ctx, question)
}

func (s *QuestionService) Update(ctx context.Context, question *model.Question) error {
	return s.QuestionRepo.Update(ctx, question)
}

func (s *QuestionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.QuestionRepo.Delete(ctx, id)
}

func (s *QuestionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Question, error) {
	return s.QuestionRepo.GetByID(ctx, id)
}

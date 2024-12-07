package service

import (
	"context"
	"golizilla/core/domain/model"
	"golizilla/core/port/repository"

	"github.com/google/uuid"
)

type IAnswerService interface {
	Create(ctx context.Context, userCtx context.Context, answer *model.Answer) (uuid.UUID, error)
	Update(ctx context.Context, userCtx context.Context, answer *model.Answer) error
	Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Answer, error)
}

type AnswerService struct {
	answerRepo repository.IAnswerRepository
}

func NewAnswerService(repo repository.IAnswerRepository) IAnswerService {
	return &AnswerService{
		answerRepo: repo,
	}
}

func (s *AnswerService) Create(ctx context.Context, userCtx context.Context, answer *model.Answer) (uuid.UUID, error) {
	return s.answerRepo.Create(ctx, userCtx, answer)
}

func (s *AnswerService) Update(ctx context.Context, userCtx context.Context, answer *model.Answer) error {
	return s.answerRepo.Update(ctx, userCtx, answer)
}

func (s *AnswerService) Delete(ctx context.Context, userCtx context.Context, id uuid.UUID) error {
	return s.answerRepo.Delete(ctx, userCtx, id)
}

func (s *AnswerService) GetByID(ctx context.Context, userCtx context.Context, id uuid.UUID) (*model.Answer, error) {
	return s.answerRepo.GetByID(ctx, userCtx, id)
}

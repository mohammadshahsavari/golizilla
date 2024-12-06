package service

import (
	"golizilla/domain/model"

	"github.com/google/uuid"
)

type ICoreService interface {
	Start(questionnaireID uuid.UUID) (*model.Question, error)
	Submit(questionID uuid.UUID, answer *model.Answer) error
	Back() (*model.Question, error)
	Next() (*model.Question, error)
	End() error
}

type CoreService struct {
}

func NewCoreService() ICoreService {
	return &CoreService{}
}

func (c *CoreService) Start(questionnaireID uuid.UUID) (*model.Question, error) {
	return nil, nil
}

func (c *CoreService) Submit(questionID uuid.UUID, answer *model.Answer) error {
	return nil
}

func (c *CoreService) Back() (*model.Question, error) {
	return nil, nil
}

func (c *CoreService) Next() (*model.Question, error) {
	return nil, nil
}

func (c *CoreService) End() error {
	return nil
}

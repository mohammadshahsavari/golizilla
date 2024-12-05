package presenter

import (
	"errors"
	"golizilla/domain/model"
	"strings"

	"github.com/google/uuid"
)

type CreateAnswerRequest struct {
	QuestionID   uuid.UUID `json:"question_id"`
	Answer       string    `json:"answer"`
	AnswerOption uint      `json:"answer_option"`
}

func (req *CreateAnswerRequest) Validate() error {
	if strings.TrimSpace(req.Answer) == "" {
		return errors.New("answer cannot be empty")
	}
	if req.QuestionID == uuid.Nil {
		return errors.New("answer cannot be empty")
	}
	return nil
}

func (req *CreateAnswerRequest) ToDomain() *model.Answer {
	return &model.Answer{
		QuestionID:   req.QuestionID,
		Answer:       req.Answer,
		AnswerOption: req.AnswerOption,
	}
}

type CreateAnswerResponse struct {
	ID uuid.UUID `json:"id"`
}

func NewCreateAnswerResponse(id uuid.UUID) CreateAnswerResponse {
	return CreateAnswerResponse{
		ID: id,
	}
}

type GetAnswerResponse struct {
	ID           uuid.UUID `json:"id"`
	QuestionID   uuid.UUID `json:"question_id"`
	Answer       string    `json:"answer"`
	AnswerOption uint      `json:"answer_option"`
}

func NewGetAnswerResponse(q *model.Answer) *GetAnswerResponse {
	return &GetAnswerResponse{
		ID:           q.ID,
		QuestionID:   q.QuestionID,
		Answer:       q.Answer,
		AnswerOption: q.AnswerOption,
	}
}

type UpdateAnswerRequest struct {
	QuestionID   uuid.UUID `json:"question_id"`
	Answer       string    `json:"answer"`
	AnswerOption uint      `json:"answer_option"`
}

func (req *UpdateAnswerRequest) Validate() error {
	if strings.TrimSpace(req.Answer) == "" {
		return errors.New("answer cannot be empty")
	}
	if req.QuestionID == uuid.Nil {
		return errors.New("answer cannot be empty")
	}
	return nil
}

func (req *UpdateAnswerRequest) ToDomain() *model.Answer {
	return &model.Answer{
		QuestionID:   req.QuestionID,
		Answer:       req.Answer,
		AnswerOption: req.AnswerOption,
	}
}

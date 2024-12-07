package presenter

import (
	"errors"
	"fmt"
	"golizilla/domain/model"
	"strings"

	"github.com/google/uuid"
)

type CreateQuestionRequest struct {
	QuestionnaireId uuid.UUID  `json:"questionnaire_id"`
	QuestionText    string     `json:"question_text"`
	Descriptive     bool       `json:"descriptive"`
	MetaDataPath    string     `json:"meta_data_path,omitempty"`
	CorrectOptionID *uuid.UUID `json:"correct_option_id,omitempty"`
	Options         []string   `json:"options,omitempty"`
}

func (req *CreateQuestionRequest) Validate() error {
	if req.QuestionnaireId == uuid.Nil {
		return errors.New("questionnaire_id cannot be empty")
	}

	if strings.TrimSpace(req.QuestionText) == "" {
		return errors.New("question text cannot be empty")
	}

	// If not descriptive, should have at least one option
	if !req.Descriptive && len(req.Options) == 0 {
		return errors.New("non-descriptive question must have at least one option")
	}
	return nil
}

func (req *CreateQuestionRequest) ToDomain() *model.Question {
	q := &model.Question{
		ID:              uuid.New(),
		QuestionnaireId: req.QuestionnaireId,
		QuestionText:    req.QuestionText,
		Descriptive:     req.Descriptive,
		MetaDataPath:    req.MetaDataPath,
		CorrectOptionID: req.CorrectOptionID,
	}

	// If it's a multiple-choice question, create options
	if !req.Descriptive && len(req.Options) > 0 {
		opts := make([]model.Option, len(req.Options))
		for i, text := range req.Options {
			opts[i] = model.Option{
				ID:   uuid.New(),
				Text: text,
				// QuestionID will be assigned when Question is created (some DB logic might handle this)
				Index: uint(i + 1),
			}
		}
		q.Options = opts
	}
	return q
}

type CreateQuestionResponse struct {
	ID uuid.UUID `json:"id"`
}

func NewCreateQuestionResponse(id uuid.UUID) CreateQuestionResponse {
	return CreateQuestionResponse{
		ID: id,
	}
}

type OptionResponse struct {
	ID    uuid.UUID `json:"id"`
	Index uint      `json:"index"`
	Text  string    `json:"text"`
}

type GetQuestionResponse struct {
	ID              uuid.UUID        `json:"id"`
	QuestionnaireId uuid.UUID        `json:"questionnaire_id"`
	Index           uint             `json:"index"`
	QuestionText    string           `json:"question_text"`
	Descriptive     bool             `json:"descriptive"`
	MetaDataPath    string           `json:"meta_data_path,omitempty"`
	CorrectOptionID *uuid.UUID       `json:"correct_option_id,omitempty"`
	Options         []OptionResponse `json:"options,omitempty"`
}

func NewGetQuestionResponse(q *model.Question) *GetQuestionResponse {
	opts := make([]OptionResponse, len(q.Options))
	for i, o := range q.Options {
		opts[i] = OptionResponse{
			ID:    o.ID,
			Index: o.Index,
			Text:  o.Text,
		}
	}

	return &GetQuestionResponse{
		ID:              q.ID,
		QuestionnaireId: q.QuestionnaireId,
		Index:           q.Index,
		QuestionText:    q.QuestionText,
		Descriptive:     q.Descriptive,
		MetaDataPath:    q.MetaDataPath,
		CorrectOptionID: q.CorrectOptionID,
		Options:         opts,
	}
}

type UpdateQuestionRequest struct {
	QuestionText    *string    `json:"question_text,omitempty"`
	Descriptive     *bool      `json:"descriptive,omitempty"`
	MetaDataPath    *string    `json:"meta_data_path,omitempty"`
	CorrectOptionID *uuid.UUID `json:"correct_option_id,omitempty"`
	Options         *[]string  `json:"options,omitempty"`
}

func (req *UpdateQuestionRequest) Validate() error {
	// If you need validation for updates, you can add here.
	// For now, it's optional fields.
	return nil
}

func (req *UpdateQuestionRequest) ToDomain(q *model.Question) *model.Question {
	if req.QuestionText != nil {
		q.QuestionText = *req.QuestionText
	}
	if req.Descriptive != nil {
		q.Descriptive = *req.Descriptive
	}
	if req.MetaDataPath != nil {
		q.MetaDataPath = *req.MetaDataPath
	}
	if req.CorrectOptionID != nil {
		q.CorrectOptionID = req.CorrectOptionID
	}

	if req.Options != nil && !q.Descriptive {
		opts := make([]model.Option, len(*req.Options))
		for i, text := range *req.Options {
			opts[i] = model.Option{
				ID:    uuid.New(),
				Text:  text,
				Index: uint(i + 1),
			}
		}
		q.Options = opts
	} else if req.Options != nil && q.Descriptive {
		// Descriptive question shouldn't have options,
		// you might want to return an error or just ignore.
		fmt.Println("Warning: Trying to set options for a descriptive question.")
	}

	return q
}

package presenter

import (
	"errors"
	"golizilla/core/domain/model"
	"strings"

	"github.com/google/uuid"
)

type CreateAnswerRequest struct {
	QuestionID   uuid.UUID  `json:"question_id"`
	Descriptive  bool       `json:"descriptive"`
	Text         *string    `json:"text,omitempty"`
	OptionID     *uuid.UUID `json:"option_id,omitempty"`
	SubmissionID uuid.UUID  `json:"submission_id"`
}

// Validate ensures the incoming request has the correct fields set.
func (req *CreateAnswerRequest) Validate() error {
	if req.QuestionID == uuid.Nil {
		return errors.New("question_id cannot be empty")
	}

	if req.SubmissionID == uuid.Nil {
		return errors.New("submission_id cannot be empty")
	}

	// If the answer is descriptive, Text should be provided
	if req.Descriptive {
		if req.Text == nil || *req.Text == "" {
			return errors.New("text answer cannot be empty for a descriptive question")
		}
		if req.OptionID != nil {
			return errors.New("option_id should be empty for a descriptive question")
		}
	} else {
		// If not descriptive, we expect an OptionID to be chosen.
		if req.OptionID == nil {
			return errors.New("option_id is required for non-descriptive answer")
		}
		if req.Text != nil {
			return errors.New("text should be empty for a descriptive question")
		}
	}
	return nil
}

// ToDomain transforms the request into a domain model Answer.
// Note: You need to supply UserID and UserSubmissionID when calling this method
// from your handler, since that context isn't in the request.
func (req *CreateAnswerRequest) ToDomain(userID uuid.UUID) *model.Answer {
	return &model.Answer{
		ID:               uuid.New(),
		QuestionID:       req.QuestionID,
		UserID:           userID,
		UserSubmissionID: req.SubmissionID,
		Descriptive:      req.Descriptive,
		Text:             req.Text,
		OptionID:         req.OptionID,
	}
}

type CreateAnswerResponse struct {
	ID uuid.UUID `json:"id"`
}

func NewCreateAnswerResponse(id uuid.UUID) CreateAnswerResponse {
	return CreateAnswerResponse{ID: id}
}

type GetAnswerResponse struct {
	ID          uuid.UUID  `json:"id"`
	QuestionID  uuid.UUID  `json:"question_id"`
	Descriptive bool       `json:"descriptive"`
	Text        *string    `json:"text,omitempty"`
	OptionID    *uuid.UUID `json:"option_id,omitempty"`
}

func NewGetAnswerResponse(a *model.Answer) *GetAnswerResponse {
	return &GetAnswerResponse{
		ID:          a.ID,
		QuestionID:  a.QuestionID,
		Descriptive: a.Descriptive,
		Text:        a.Text,
		OptionID:    a.OptionID,
	}
}

type UpdateAnswerRequest struct {
	Descriptive bool       `json:"descriptive"`
	Text        *string    `json:"text,omitempty"`
	OptionID    *uuid.UUID `json:"option_id,omitempty"`
}

func (req *UpdateAnswerRequest) Validate() error {
	if req.Descriptive && (req.Text == nil || strings.TrimSpace(*req.Text) == "") {
		return errors.New("text answer cannot be empty for descriptive")
	}
	if !req.Descriptive && req.OptionID == nil {
		return errors.New("option_id is required for non-descriptive answer")
	}

	if req.Descriptive {
		if req.OptionID != nil {
			return errors.New("option_id should be empty for a descriptive question")
		}
	} else {
		if req.Text != nil {
			return errors.New("text should be empty for a descriptive question")
		}
	}
	return nil
}

func (req *UpdateAnswerRequest) ToDomain(a *model.Answer) *model.Answer {
	a.Descriptive = req.Descriptive
	a.Text = req.Text
	a.OptionID = req.OptionID
	return a
}

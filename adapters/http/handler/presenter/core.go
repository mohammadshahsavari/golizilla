package presenter

import (
	"errors"
	"golizilla/core/domain/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// StartRequest represents the request structure for starting a questionnaire.
type StartRequest struct {
	QuestionnaireID uuid.UUID
	UserID          uuid.UUID
}

func (r *StartRequest) ParseAndValidate(c *fiber.Ctx) error {
	qID, err := uuid.Parse(c.Params("questionnaire_id"))
	if err != nil {
		return errors.New("invalid questionnaire_id format")
	}
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return errors.New("user_id not found or invalid")
	}
	r.QuestionnaireID = qID
	r.UserID = userID
	return nil
}

// SubmitRequest represents the request structure for submitting an answer.
type SubmitRequest struct {
	SubmissionID uuid.UUID  `json:"submission_id"`
	QuestionID   uuid.UUID  `json:"question_id"`
	UserID       uuid.UUID  `json:"-"` // Set dynamically
	Descriptive  bool       `json:"descriptive"`
	Text         *string    `json:"text,omitempty"`
	OptionID     *uuid.UUID `json:"option_id,omitempty"`
}

func (r *SubmitRequest) ParseAndValidate(c *fiber.Ctx) error {
	if err := c.BodyParser(r); err != nil {
		return errors.New("invalid request format")
	}

	// Validate required fields
	if r.SubmissionID == uuid.Nil {
		return errors.New("submission_id is required")
	}
	if r.QuestionID == uuid.Nil {
		return errors.New("question_id is required")
	}
	if !r.Descriptive && r.OptionID == nil {
		return errors.New("option_id is required for non-descriptive answers")
	}
	if r.Descriptive && (r.Text == nil || *r.Text == "") {
		return errors.New("text is required for descriptive answers")
	}

	return nil
}

func (r *SubmitRequest) ToDomain() *model.Answer {
	return &model.Answer{
		QuestionID:       r.QuestionID,
		UserID:           r.UserID,
		UserSubmissionID: r.SubmissionID,
		Descriptive:      r.Descriptive,
		Text:             r.Text,
		OptionID:         r.OptionID,
	}
}

// NavigationRequest represents a request to navigate within a submission.
type NavigationRequest struct {
	SubmissionID uuid.UUID
}

func (r *NavigationRequest) ParseAndValidate(c *fiber.Ctx) error {
	var req struct {
		SubmissionID string `json:"submission_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return errors.New("invalid request body")
	}

	sID, err := uuid.Parse(req.SubmissionID)
	if err != nil {
		return errors.New("invalid submission_id format")
	}

	r.SubmissionID = sID
	return nil
}

// EndRequest represents a request to end a submission.
type EndRequest struct {
	SubmissionID uuid.UUID
}

func (r *EndRequest) ParseAndValidate(c *fiber.Ctx) error {
	var req struct {
		SubmissionID string `json:"submission_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return errors.New("invalid request body")
	}

	sID, err := uuid.Parse(req.SubmissionID)
	if err != nil {
		return errors.New("invalid submission_id format")
	}

	r.SubmissionID = sID
	return nil
}

// StartResponse prepares the response structure for starting a questionnaire.
func NewStartResponse(submissionID uuid.UUID, question *model.Question) map[string]interface{} {
	return map[string]interface{}{
		"submission_id":    submissionID,
		"id":               question.ID,
		"questionnaire_id": question.QuestionnaireId,
		"index":            question.Index,
		"question_text":    question.QuestionText,
		"descriptive":      question.Descriptive,
		"options":          question.Options,
	}
}

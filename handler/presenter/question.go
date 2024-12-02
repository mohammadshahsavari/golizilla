package presenter

import (
	"errors"
	"fmt"
	"golizilla/domain/model"
	"strings"

	"github.com/google/uuid"
)

type CreateQuestionRequest struct {
	QuestionnaireId uuid.UUID `json:"questionId"`
	QuestionText    string    `json:"questionText"`
	Descriptive     bool      `json:"descriptive"`
	OptionsCount    uint      `json:"optionsCount"`
	CorrectOption   uint      `json:"correctOption"`
	MetaDataPath    string    `json:"metaDataPath,omitempty"`
	OptionsText     string    `json:"optionsText"` // string or []byte ?
	// Answers         []*Answer `json:""`
}

func (req *CreateQuestionRequest) Validate() error {
	// Validate QuestionnaireId
	if req.QuestionnaireId == uuid.Nil {
		return errors.New("questionnaireId cannot be empty")
	}

	// Validate QuestionText
	if strings.TrimSpace(req.QuestionText) == "" {
		return errors.New("questionText cannot be empty")
	}

	// Validate CorrectOption
	if req.CorrectOption > req.OptionsCount {
		return fmt.Errorf("correctOption must be less or equal than %d", req.OptionsCount)
	}

	// // Validate OptionsText
	// if len(req.OptionsText) != int(req.OptionsCount) {
	// 	return fmt.Errorf("optionsText must contain exactly %d options", req.OptionsCount)
	// }

	// if _, err := json.Marshal(req.OptionsText); err != nil {
	// 	return fmt.Errorf("Error marshaling options:", err)
	// }

	// for _, option := range req.OptionsText {
	// 	if option == "" {
	// 		return errors.New("each option in optionsText must be non-empty")
	// 	}
	// }

	return nil
}

func (req *CreateQuestionRequest) ToDomain() *model.Question {
	return &model.Question{
		QuestionnaireId: req.QuestionnaireId,
		QuestionText:    req.QuestionText,
		Descriptive:     req.Descriptive,
		OptionsCount:    req.OptionsCount,
		CorrectOption:   req.CorrectOption,
		MetaDataPath:    req.MetaDataPath,
		OptionsText:     req.OptionsText,
	}
}

type CreateQuestionResponse struct {
	ID uuid.UUID `json:"id"`
}

func NewCreateQuestionResponse(id uuid.UUID) CreateQuestionResponse {
	return CreateQuestionResponse{
		ID: id,
	}
}

type GetQuestionResponse struct {
	ID              uuid.UUID `json:"id"`
	QuestionnaireId uuid.UUID `json:"questionnaireId"`
	Index           uint      `json:"index"`
	QuestionText    string    `json:"questionText"`
	Descriptive     bool      `json:"descriptive"`
	OptionsCount    uint      `json:"optionsCount"`
	CorrectOption   uint      `json:"correctOption"`
	MetaDataPath    string    `json:"metaDataPath,omitempty"`
	OptionsText     string    `json:"optionsText"`
	// Answers         []*Answer
}

func NewGetQuestionResponse(q *model.Question) *GetQuestionResponse {
	return &GetQuestionResponse{
		ID:              q.ID,
		QuestionnaireId: q.QuestionnaireId,
		Index:           q.Index,
		QuestionText:    q.QuestionText,
		Descriptive:     q.Descriptive,
		OptionsCount:    q.OptionsCount,
		CorrectOption:   q.CorrectOption,
		MetaDataPath:    q.MetaDataPath,
		OptionsText:     q.OptionsText,
	}
}

type UpdateQuestionRequest struct {
	QuestionText  string `json:"questionText"`
	Descriptive   bool   `json:"descriptive"`
	OptionsCount  uint   `json:"optionsCount"`
	CorrectOption uint   `json:"correctOption"`
	MetaDataPath  string `json:"metaDataPath"`
	OptionsText   string `json:"optionsText"`
}

func (req *UpdateQuestionRequest) Validate() error {
	// Validate QuestionText
	if strings.TrimSpace(req.QuestionText) == "" {
		return errors.New("questionText cannot be empty")
	}

	// Validate CorrectOption
	if req.CorrectOption > req.OptionsCount {
		return fmt.Errorf("correctOption must be less or equal than %d", req.OptionsCount)
	}

	// // Validate OptionsText
	// if len(req.OptionsText) != int(req.OptionsCount) {
	// 	return fmt.Errorf("optionsText must contain exactly %d options", req.OptionsCount)
	// }

	// if _, err := json.Marshal(req.OptionsText); err != nil {
	// 	return fmt.Errorf("Error marshaling options:", err)
	// }

	// for _, option := range req.OptionsText {
	// 	if option == "" {
	// 		return errors.New("each option in optionsText must be non-empty")
	// 	}
	// }

	return nil
}

func (req *UpdateQuestionRequest) ToDomain() *model.Question {
	return &model.Question{
		ID:              [16]byte{},
		QuestionnaireId: [16]byte{},
		Index:           0,
		QuestionText:    req.QuestionText,
		Descriptive:     req.Descriptive,
		OptionsCount:    req.OptionsCount,
		CorrectOption:   req.CorrectOption,
		MetaDataPath:    req.MetaDataPath,
		OptionsText:     req.OptionsText,
		SelectedOption:  0,
		Answers:         []*model.Answer{},
	}
}

package presenter

import (
	"errors"
	"fmt"
	"golizilla/domain/model"
	"strings"

	"github.com/google/uuid"
)

type CreateQuestionRequest struct {
	QuestionnaireId uuid.UUID `json:"question_id"`
	QuestionText    string    `json:"question_text"`
	Descriptive     bool      `json:"descriptive"`
	OptionsCount    uint      `json:"options_count"`
	CorrectOption   uint      `json:"correct_option"`
	MetaDataPath    string    `json:"meta_data_path,omitempty"`
	OptionsText     string    `json:"options_text"` // string or []byte ?
	// Answers         []*Answer `json:""`
}

func (req *CreateQuestionRequest) Validate() error {
	// Validate QuestionnaireId
	if req.QuestionnaireId == uuid.Nil {
		return errors.New("questionnaire id cannot be empty")
	}

	// Validate QuestionText
	if strings.TrimSpace(req.QuestionText) == "" {
		return errors.New("question text cannot be empty")
	}

	// Validate CorrectOption
	if req.CorrectOption > req.OptionsCount {
		return fmt.Errorf("correct option must be less or equal than %d", req.OptionsCount)
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
	QuestionnaireId uuid.UUID `json:"questionnaire_id"`
	Index           uint      `json:"index"`
	QuestionText    string    `json:"question_text"`
	Descriptive     bool      `json:"descriptive"`
	OptionsCount    uint      `json:"options_count"`
	CorrectOption   uint      `json:"correct_option"`
	MetaDataPath    string    `json:"meta_data_path,omitempty"`
	OptionsText     string    `json:"options_text"`
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
	QuestionText  string `json:"question_text"`
	Descriptive   bool   `json:"descriptive"`
	OptionsCount  uint   `json:"options_count"`
	CorrectOption uint   `json:"correct_option"`
	MetaDataPath  string `json:"meta_data_path"`
	OptionsText   string `json:"options_text"`
}

func (req *UpdateQuestionRequest) Validate() error {
	// Validate QuestionText
	if strings.TrimSpace(req.QuestionText) == "" {
		return errors.New("question text cannot be empty")
	}

	// Validate CorrectOption
	if req.CorrectOption > req.OptionsCount {
		return fmt.Errorf("correct option must be less or equal than %d", req.OptionsCount)
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

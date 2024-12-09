package presenter

import (
	"errors"
	"golizilla/core/domain/model"
	"time"

	"github.com/google/uuid"
)

type CreateQuestionnaireRequest struct {
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	Random         bool          `json:"random"`
	BackCompatible bool          `json:"back_compatible"`
	Title          string        `json:"title"`
	AnswerTime     time.Duration `json:"answer_time"`
	Anonymous      bool          `json:"anonymous"`
	SubmitLimit    uint          `json:"submit_limit,omitempty"`
	//TODO: Questions
}

type GiveAcessRequest struct {
	UserIDs    []uuid.UUID `json:"user_ids"`
	Privileges []string    `json:"privileges"`
	AllUsers   bool        `json:"all_users"`
}

type DeleteAcessRequest struct {
	UserIDs    []uuid.UUID `json:"user_ids"`
	Privileges []string    `json:"privileges"`
	AllUsers   bool        `json:"all_users"`
}

type UpdateQuestionnaireRequest struct {
	ID             uuid.UUID      `json:"id"` // Mandatory for updates
	Random         *bool          `json:"random,omitempty"`
	BackCompatible *bool          `json:"back_compatible,omitempty"`
	Title          *string        `json:"title,omitempty"`
	AnswerTime     *time.Duration `json:"answer_time,omitempty"`
	Anonymous      *bool          `json:"anonymous,omitempty"`
}

type CreateQuestionnaireResponseData struct {
	Id uuid.UUID `json:"id"`
}

type GetQuestionnaireResponseData struct {
	Id                 uuid.UUID     `json:"id"`
	OwnerId            uuid.UUID     `json:"owner_id"`
	CreatedTime        time.Time     `json:"created_time"`
	StartTime          time.Time     `json:"start_time"`
	EndTime            time.Time     `json:"end_time"`
	Random             bool          `json:"random"`
	BackCompatible     bool          `json:"back_compatible"`
	Title              string        `json:"title"`
	AnswerTime         time.Duration `json:"answer_time"`
	ParticipationCount uint          `json:"particpation_count"`
	Anonymous          bool          `json:"anonymous"`
}

func (req *CreateQuestionnaireRequest) Validate() error {
	// Validate Title
	if req.Title == "" {
		return errors.New("title can't be empty")
	}

	// Validate AnswerTime (should be positive)
	if req.AnswerTime <= 0 {
		return errors.New("answer time must be greater than zero")
	}

	if req.StartTime.After(req.EndTime) {
		return errors.New("start time cannot be after end time")
	}
	if req.EndTime.Before(time.Now()) {
		return errors.New("end time must be in the future")
	}

	return nil
}

func (req *CreateQuestionnaireRequest) ToDomain() *model.Questionnaire {
	return &model.Questionnaire{
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		CreatedTime:    time.Now(),
		Random:         req.Random,
		BackCompatible: req.BackCompatible,
		Title:          req.Title,
		AnswerTime:     req.AnswerTime,
		Anonymous:      req.Anonymous,
	}
}

func (r *UpdateQuestionnaireRequest) Validate() error {
	if r.ID == uuid.Nil {
		return errors.New("ID is required for updating a questionnaire")
	}
	return nil
}

func (r *UpdateQuestionnaireRequest) ToDomain() map[string]interface{} {
	updateFields := map[string]interface{}{}

	if r.Random != nil {
		updateFields["random"] = *r.Random
	}
	if r.BackCompatible != nil {
		updateFields["back_compatible"] = *r.BackCompatible
	}
	if r.Title != nil {
		updateFields["title"] = *r.Title
	}
	if r.AnswerTime != nil {
		updateFields["answer_time"] = *r.AnswerTime
	}
	if r.Anonymous != nil {
		updateFields["anonymous"] = *r.Anonymous
	}

	return updateFields
}

func NewCreateQuestionnaireResponse(id uuid.UUID) Response {
	return Response{
		Success: true,
		Data: CreateQuestionnaireResponseData{
			Id: id,
		},
	}
}

func NewGetQuestionnaireResponse(data *model.Questionnaire) Response {
	return Response{
		Success: true,
		Data: GetQuestionnaireResponseData{
			Id:                 data.Id,
			OwnerId:            data.OwnerId,
			CreatedTime:        data.CreatedTime,
			StartTime:          data.StartTime,
			EndTime:            data.EndTime,
			Random:             data.Random,
			BackCompatible:     data.BackCompatible,
			Title:              data.Title,
			AnswerTime:         data.AnswerTime,
			ParticipationCount: data.ParticipationCount,
			Anonymous:          data.Anonymous,
		},
	}
}

func NewGetQuestionnairesResponse(data []model.Questionnaire) Response {
	var resultData []GetQuestionnaireResponseData

	for _, item := range data {
		resultData = append(resultData, GetQuestionnaireResponseData{
			Id:                 item.Id,
			OwnerId:            item.OwnerId,
			CreatedTime:        item.CreatedTime,
			StartTime:          item.StartTime,
			EndTime:            item.EndTime,
			Random:             item.Random,
			BackCompatible:     item.BackCompatible,
			Title:              item.Title,
			AnswerTime:         item.AnswerTime,
			ParticipationCount: item.ParticipationCount,
			Anonymous:          item.Anonymous,
		})
	}
	return Response{
		Success: true,
		Data:    resultData,
	}
}

package presenter

import (
	"errors"
	"golizilla/domain/model"
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
	//TODO: Questions
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
	if req.EndTime.Before(req.StartTime) || req.EndTime.Before(time.Now()) {
		return errors.New("start and end time is not valid")
	}
	if req.Title == "" {
		return errors.New("title can't be empty")
	}

	return nil
}

func (req *CreateQuestionnaireRequest) ToDomain() *model.Questionnaire {
	return &model.Questionnaire{
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Random:         req.Random,
		BackCompatible: req.BackCompatible,
		Title:          req.Title,
		AnswerTime:     req.AnswerTime,
		Anonymous:      req.Anonymous,
	}
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

func NewGetQuestionnaireesResponse(data []model.Questionnaire) Response {
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

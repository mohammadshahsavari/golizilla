package presenter

import (
	"errors"
	"golizilla/domain/model"
	"time"

	"github.com/google/uuid"
)

type CreateQuestionnaryRequest struct {
	StartTime      time.Time     `json:"startTime"`
	EndTime        time.Time     `json:"endTime"`
	Random         bool          `json:"random"`
	BackCompatible bool          `json:"backCompatible"`
	Title          string        `json:"title"`
	AnswerTime     time.Duration `json:"answerTime"`
	//TODO: Questions
}

type CreateQuestionnaryResponseData struct {
	Id uuid.UUID `json:"id"`
}

type GetQuestionnaryResponseData struct {
	Id                uuid.UUID     `json:"id"`
	OwnerId           uuid.UUID     `json:"ownerId"`
	CreatedTime       time.Time     `json:"createdTime"`
	StartTime         time.Time     `json:"startTime"`
	EndTime           time.Time     `json:"endTime"`
	Random            bool          `json:"random"`
	BackCompatible    bool          `json:"backCompatible"`
	Title             string        `json:"title"`
	AnswerTime        time.Duration `json:"answerTime"`
	ParticpationCount uint          `json:"particpationCount"`
}

func (req *CreateQuestionnaryRequest) Validate() error {
	if req.EndTime.Before(req.StartTime) || req.EndTime.Before(time.Now()) {
		return errors.New("start and end time is not valid")
	}
	if req.Title == "" {
		return errors.New("title can't be empty")
	}

	return nil
}

func (req *CreateQuestionnaryRequest) ToDomain() *model.Questionnaire {
	return &model.Questionnaire{
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Random:         req.Random,
		BackCompatible: req.BackCompatible,
		Title:          req.Title,
		AnswerTime:     req.AnswerTime,
	}
}

func NewCreateQuestionnaryResponse(id uuid.UUID) Response {
	return Response{
		Success: true,
		Data: CreateQuestionnaryResponseData{
			Id: id,
		},
	}
}

func NewGetQuestionnaryResponse(data *model.Questionnaire) Response {
	return Response{
		Success: true,
		Data: GetQuestionnaryResponseData{
			Id:                data.Id,
			OwnerId:           data.OwnerId,
			CreatedTime:       data.CreatedTime,
			StartTime:         data.StartTime,
			EndTime:           data.EndTime,
			Random:            data.Random,
			BackCompatible:    data.BackCompatible,
			Title:             data.Title,
			AnswerTime:        data.AnswerTime,
			ParticpationCount: data.ParticpationCount,
		},
	}
}

func NewGetQuestionnariesResponse(data []model.Questionnaire) Response {
	var resultData []GetQuestionnaryResponseData

	for _, item := range data {
		resultData = append(resultData, GetQuestionnaryResponseData{
			Id:                item.Id,
			OwnerId:           item.OwnerId,
			CreatedTime:       item.CreatedTime,
			StartTime:         item.StartTime,
			EndTime:           item.EndTime,
			Random:            item.Random,
			BackCompatible:    item.BackCompatible,
			Title:             item.Title,
			AnswerTime:        item.AnswerTime,
			ParticpationCount: item.ParticpationCount,
		})
	}
	return Response{
		Success: true,
		Data:    resultData,
	}
}

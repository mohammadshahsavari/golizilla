package handler

import (
	"context"
	"fmt"
	"golizilla/handler/presenter"
	"golizilla/internal/apperrors"
	"golizilla/service"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionnaireHandler struct {
	questionnaireService service.IQuestionnaireService
}

func NewQuestionnaireHandler(questionnaireService service.IQuestionnaireService) *QuestionnaireHandler {
	return &QuestionnaireHandler{
		questionnaireService: questionnaireService,
	}
}

func (q *QuestionnaireHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()
	var request presenter.CreateQuestionnaireRequest
	if err := c.BodyParser(&request); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	if err := request.Validate(); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	userModel := request.ToDomain()
	userModel.OwnerId = c.Locals("user_id").(uuid.UUID)
	if id, err := q.questionnaireService.Create(ctx, userModel); err != nil {
		//log
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	} else {
		return presenter.Send(c, fiber.StatusOK, true, "Questionnaire created successfully", presenter.NewCreateQuestionnaireResponse(id), nil)
	}
}

func (q *QuestionnaireHandler) Delete(c *fiber.Ctx) error {
	ctx := c.Context()
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid ID format")
	}
	err = q.questionnaireService.Delete(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presenter.SendError(c, fiber.StatusNotFound, err.Error())
		}

		//log
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "Deleted", nil, nil)
}

func (q *QuestionnaireHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		// log
		fmt.Printf("[Update Questionnaire] invalid input error: %v", err)
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	var request presenter.UpdateQuestionnaireRequest
	request.ID = id
	if err := c.BodyParser(&request); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	if err := request.Validate(); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Map fields to update
	updateFields := request.ToDomain()

	// Pass the ID and update fields to the service
	if err := q.questionnaireService.Update(ctx, request.ID, updateFields); err != nil {
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "Questionnaire updated successfully", nil, nil)
}

func (q *QuestionnaireHandler) GetById(c *fiber.Ctx) error {
	ctx := c.Context()
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid ID format")
	}
	questionnaire, err := q.questionnaireService.GetById(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presenter.SendError(c, fiber.StatusNotFound, err.Error())
		}

		//log
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "", presenter.NewGetQuestionnaireResponse(questionnaire), nil)
}

func (q *QuestionnaireHandler) GetByOwnerId(c *fiber.Ctx) error {
	ctx := c.Context()
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid ID format")
	}
	questionnaires, err := q.questionnaireService.GetByOwnerId(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presenter.SendError(c, fiber.StatusNotFound, err.Error())
		}

		//log
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "", presenter.NewGetQuestionnairesResponse(questionnaires), nil)
}

func (q *QuestionnaireHandler) GetResults(c *websocket.Conn) {
	idString := c.Params("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
		return
	}
	_, err = q.questionnaireService.GetById(context.Background(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
			return
		}
		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
		return
	}

	var lastValue uint = 0
	for {
		questionnaire, err := q.questionnaireService.GetById(context.Background(), id)
		if err != nil {
			//log
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
			break
		}
		if lastValue != questionnaire.ParticipationCount {
			lastValue = questionnaire.ParticipationCount
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", lastValue)))
		}
		time.Sleep(time.Second * 10)
	}
}

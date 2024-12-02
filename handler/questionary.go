package handler

import (
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

type QuestionnaryHandler struct {
	questionnaryService service.IQuestionnaireService
}

func NewQuestionnaryHandler(questionnaryService service.IQuestionnaireService) *QuestionnaryHandler {
	return &QuestionnaryHandler{
		questionnaryService: questionnaryService,
	}
}

func (q *QuestionnaryHandler) Create(c *fiber.Ctx) error {
	var request presenter.CreateQuestionnaryRequest
	if err := c.BodyParser(&request); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	if err := request.Validate(); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	userModel := request.ToDomain()
	userModel.OwnerId = c.Locals("user_id").(uuid.UUID)
	if id, err := q.questionnaryService.Create(userModel); err != nil {
		//log
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	} else {
		return presenter.Send(c, fiber.StatusOK, true, "Questionnary created successfully", presenter.NewCreateQuestionnaryResponse(id), nil)
	}
}

func (q *QuestionnaryHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid ID format")
	}
	err = q.questionnaryService.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presenter.SendError(c, fiber.StatusNotFound, err.Error())
		}

		//log
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "Deleted", nil, nil)
}

func (q *QuestionnaryHandler) Update(c *fiber.Ctx) error {
	var request presenter.CreateQuestionnaryRequest
	if err := c.BodyParser(&request); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	if err := request.Validate(); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	userModel := request.ToDomain()

	if err := q.questionnaryService.Update(userModel); err != nil {
		//log
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "Updated", nil, nil)
}

func (q *QuestionnaryHandler) GetById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid ID format")
	}
	questionary, err := q.questionnaryService.GetById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presenter.SendError(c, fiber.StatusNotFound, err.Error())
		}

		//log
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "", presenter.NewGetQuestionnaryResponse(questionary), nil)
}

func (q *QuestionnaryHandler) GetByOwnerId(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid ID format")
	}
	questionaries, err := q.questionnaryService.GetByOwnerId(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presenter.SendError(c, fiber.StatusNotFound, err.Error())
		}

		//log
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "", presenter.NewGetQuestionnariesResponse(questionaries), nil)
}

func (q *QuestionnaryHandler) GetResults(c *websocket.Conn) {
	idString := c.Params("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
		return
	}
	_, err = q.questionnaryService.GetById(id)
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
		questionary, err := q.questionnaryService.GetById(id)
		if err != nil {
			//log
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
			break
		}
		if lastValue != questionary.ParticpationCount {
			lastValue = questionary.ParticpationCount
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", lastValue)))
		}
		time.Sleep(time.Second * 10)
	}
}

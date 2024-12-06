package handler

import (
	"context"
	"fmt"
	"golizilla/handler/presenter"
	"golizilla/internal/apperrors"
	"golizilla/internal/logmessages"
	"golizilla/persistence/logger"
	"golizilla/service"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionnariHandler struct {
	questionnariService service.IQuestionnaireService
}

func NewQuestionnariHandler(questionnariService service.IQuestionnaireService) *QuestionnariHandler {
	return &QuestionnariHandler{
		questionnariService: questionnariService,
	}
}

func (q *QuestionnariHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireCreateBegin,
	})

	var request presenter.CreateQuestionnariRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	userModel := request.ToDomain()
	userModel.OwnerId = c.Locals("user_id").(uuid.UUID)
	id, err := q.questionnariService.Create(ctx, userModel);
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	err = presenter.Send(c, fiber.StatusOK, true, "Questionnari created successfully", presenter.NewCreateQuestionnariResponse(id), nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireCreateSuccessful,
	})
	
	return nil
}

func (q *QuestionnariHandler) Delete(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireDeleteBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid ID format")
	}
	err = q.questionnariService.Delete(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			return presenter.SendError(c, fiber.StatusNotFound, err.Error())
		}
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	err = presenter.Send(c, fiber.StatusOK, true, "Deleted", nil, nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireDeleteSuccessful,
	})

	return nil
}

func (q *QuestionnariHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireUpdateBegin,
	})

	var request presenter.CreateQuestionnariRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	userModel := request.ToDomain()

	if err := q.questionnariService.Update(ctx, userModel); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	err := presenter.Send(c, fiber.StatusOK, true, "Updated", nil, nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireUpdateSuccessful,
	})

	return nil
}

func (q *QuestionnariHandler) GetById(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireGetByIdBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid ID format")
	}

	questionnari, err := q.questionnariService.GetById(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			return presenter.SendError(c, fiber.StatusNotFound, err.Error())
		}

		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	err = presenter.Send(c, fiber.StatusOK, true, "", presenter.NewGetQuestionnariResponse(questionnari), nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireGetByIdSuccessful,
	})

	return nil
}

func (q *QuestionnariHandler) GetByOwnerId(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireGetByOwnerIdBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid ID format")
	}

	questionnaries, err := q.questionnariService.GetByOwnerId(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			return presenter.SendError(c, fiber.StatusNotFound, err.Error())
		}

		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	err = presenter.Send(c, fiber.StatusOK, true, "", presenter.NewGetQuestionnariesResponse(questionnaries), nil)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return err
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireGetByOwnerIdSuccessful,
	})

	return nil
}

func (q *QuestionnariHandler) GetResults(c *websocket.Conn) {
	ctx := context.Background()
	
	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireGetResultsBegin,
	})

	idString := c.Params("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
		return
	}
	_, err = q.questionnariService.GetById(context.Background(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
			return
		}

		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
		return
	}

	var lastValue uint = 0
	for {
		questionnari, err := q.questionnariService.GetById(context.Background(), id)
		if err != nil {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
			break
		}
		if lastValue != questionnari.ParticipationCount {
			lastValue = questionnari.ParticipationCount
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", lastValue)))
		}
		time.Sleep(time.Second * 10)
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireGetResultsEnd,
	})
}

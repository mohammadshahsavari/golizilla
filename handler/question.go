package handler

import (
	"errors"
	"fmt"
	"golizilla/handler/presenter"
	"golizilla/internal/apperrors"
	"golizilla/internal/logmessages"
	"golizilla/persistence/logger"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionHandler struct {
	QuestionService service.IQuestionService
}

func NewQuestionHandler(questionService service.IQuestionService) *QuestionHandler {
	return &QuestionHandler{
		QuestionService: questionService,
	}
}

func (h *QuestionHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionHandler,
		Message: logmessages.LogQuestionCreateBegin,
	})

	var request presenter.CreateQuestionRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	question := request.ToDomain()

	id, err := h.QuestionService.Create(ctx, c.UserContext(), question)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusInternalServerError,
			apperrors.ErrInternalServerError.Error(),
		)
	}

	return presenter.Send(c,
		fiber.StatusOK, true,
		"question successfully created",
		presenter.NewCreateQuestionResponse(id),
		nil,
	)
}

func (h *QuestionHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionHandler,
		Message: logmessages.LogQuestionUpdateBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	var request presenter.UpdateQuestionRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	question, err := h.QuestionService.GetByID(ctx, c.UserContext(), id)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c,
				fiber.StatusNotFound,
				apperrors.ErrNotFound.Error(),
			)
		}
		return presenter.SendError(c,
			fiber.StatusInternalServerError,
			apperrors.ErrInternalServerError.Error(),
		)
	}

	updatedQuestion := request.ToDomain(question)
	err = h.QuestionService.Update(ctx, c.UserContext(), updatedQuestion)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusInternalServerError,
			apperrors.ErrInternalServerError.Error(),
		)
	}

	return presenter.Send(c,
		fiber.StatusOK, true,
		"question successfully updated",
		nil,
		nil,
	)
}

func (h *QuestionHandler) Delete(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionHandler,
		Message: logmessages.LogQuestionDeleteBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	err = h.QuestionService.Delete(ctx, c.UserContext(), id)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c,
				fiber.StatusNotFound,
				apperrors.ErrNotFound.Error(),
			)
		}
		return presenter.SendError(c,
			fiber.StatusInternalServerError,
			apperrors.ErrInternalServerError.Error(),
		)
	}

	return presenter.Send(c,
		fiber.StatusOK,
		true,
		fmt.Sprintf("question with id: (%v) successfully deleted", id),
		nil,
		nil,
	)
}

func (h *QuestionHandler) GetByID(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionHandler,
		Message: logmessages.LogQuestionGetByIDBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	question, err := h.QuestionService.GetByID(ctx, c.UserContext(), id)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionHandler,
			Message: err.Error(),
		})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c,
				fiber.StatusNotFound,
				apperrors.ErrNotFound.Error(),
			)
		}
		return presenter.SendError(c,
			fiber.StatusInternalServerError,
			apperrors.ErrInternalServerError.Error(),
		)
	}

	return presenter.Send(c,
		fiber.StatusOK,
		true,
		"question successfully fetched",
		presenter.NewGetQuestionResponse(question),
		nil,
	)
}

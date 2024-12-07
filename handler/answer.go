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

type AnswerHandler struct {
	answerService service.IAnswerService
}

func NewAnswerHandler(answerService service.IAnswerService) *AnswerHandler {
	return &AnswerHandler{
		answerService: answerService,
	}
}

func (h *AnswerHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogAnswerHandler,
		Message: logmessages.LogAnswerCreateBegin,
	})

	var request presenter.CreateAnswerRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
	}

	// Suppose submission_id is also given in request or from c.Locals
	// For a standalone answer, you must know which submission this answer belongs to.
	submissionID := uuid.Nil // obtain from request or context if needed

	answerDomain := request.ToDomain(userID, submissionID)
	id, err := h.answerService.Create(c.Context(), c.UserContext(), answerDomain)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusInternalServerError,
			apperrors.ErrInternalServerError.Error(),
		)
	}

	return presenter.Send(c,
		fiber.StatusOK, true,
		"Answer successfully created",
		presenter.NewCreateAnswerResponse(id),
		nil,
	)
}

func (h *AnswerHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogAnswerHandler,
		Message: logmessages.LogAnswerUpdateBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	var request presenter.UpdateAnswerRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	if err := request.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	answer, err := h.answerService.GetByID(ctx, c.UserContext(), id)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
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

	updatedAnswer := request.ToDomain(answer)
	err = h.answerService.Update(ctx, c.UserContext(), updatedAnswer)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusInternalServerError,
			apperrors.ErrInternalServerError.Error(),
		)
	}

	return presenter.Send(c,
		fiber.StatusOK, true,
		"Answer successfully updated",
		nil,
		nil,
	)
}

func (h *AnswerHandler) Delete(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogAnswerHandler,
		Message: logmessages.LogAnswerDeleteBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	err = h.answerService.Delete(ctx, c.UserContext(), id)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
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
		fmt.Sprintf("Answer with id: (%v) successfully deleted", id),
		nil,
		nil,
	)
}

func (h *AnswerHandler) GetByID(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogAnswerHandler,
		Message: logmessages.LogAnswerGetByIDBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	answer, err := h.answerService.GetByID(ctx, c.UserContext(), id)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
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
		"Answer successfully fetched",
		presenter.NewGetAnswerResponse(answer),
		nil,
	)
}

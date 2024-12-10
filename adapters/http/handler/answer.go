package handler

import (
	"errors"
	"fmt"
	"golizilla/adapters/http/handler/presenter"
	"golizilla/adapters/persistence/logger"
	"golizilla/core/service"
	"golizilla/internal/apperrors"
	logmessages "golizilla/internal/logmessages"
	privilegeconstants "golizilla/internal/privilege"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerHandler struct {
	answerService       service.IAnswerService
	qustionService      service.IQuestionService
	questionnariService service.IQuestionnaireService
	roleService         service.IRoleService
}

func NewAnswerHandler(answerService service.IAnswerService,
	qustionService service.IQuestionService,
	questionnariService service.IQuestionnaireService,
	roleService service.IRoleService) *AnswerHandler {
	return &AnswerHandler{
		answerService:       answerService,
		qustionService:      qustionService,
		questionnariService: questionnariService,
		roleService:         roleService,
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
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAnswerHandler,
			Message: logmessages.LogCastUserIdError,
		})
		return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
	}
	if answer.UserSubmissionID != userID {
		question, err := h.qustionService.GetByID(ctx, c.UserContext(), answer.QuestionID)
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

		questionnari, err := h.questionnariService.GetById(ctx, c.UserContext(), question.QuestionnaireId)
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
		if questionnari.OwnerId != userID {
			hasPrvilege, err := h.roleService.HasPrivilegesOnInsance(ctx, c.UserContext(), userID, questionnari.Id, privilegeconstants.SeeVoteOnInstance)
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

			if !hasPrvilege {
				logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
					Service: logmessages.LogQuestionnaireHandler,
					Message: logmessages.LogLackOfAuthorization,
				})
				return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrLackOfAuthorization.Error())
			}
		}
	}

	return presenter.Send(c,
		fiber.StatusOK,
		true,
		"Answer successfully fetched",
		presenter.NewGetAnswerResponse(answer),
		nil,
	)
}

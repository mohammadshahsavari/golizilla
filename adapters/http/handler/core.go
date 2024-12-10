package handler

import (
	"errors"
	"golizilla/adapters/http/handler/presenter"
	"golizilla/adapters/persistence/logger"
	"golizilla/core/service"
	"golizilla/internal/apperrors"
	"golizilla/internal/logmessages"
	privilegeconstants "golizilla/internal/privilege"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CoreHandler struct {
	coreService          service.ICoreService
	roleService          service.IRoleService
	questionnaireService service.IQuestionnaireService
}

func NewCoreHandler(coreService service.ICoreService, roleService service.IRoleService, questionnaireService service.IQuestionnaireService) *CoreHandler {
	return &CoreHandler{
		coreService:          coreService,
		roleService:          roleService,
		questionnaireService: questionnaireService,
	}
}

// StartHandler initializes a questionnaire session for the user and returns the first question.
func (h *CoreHandler) StartHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: "starting questionnaire",
	})

	// Parse and validate request
	req := &presenter.StartRequest{}
	if err := req.ParseAndValidate(c); err != nil {
		return presenter.SendError(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	// Check if the questionnaire is active
	isAnonymous, err := h.questionnaireService.IsQuestionnaireAnonymous(ctx, c.UserContext(), req.QuestionnaireID)
	if err != nil {
		if errors.Is(err, apperrors.ErrQuestionnaireNotFound) {
			logger.GetLogger().LogWarningFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: apperrors.ErrQuestionsNotFound.Error(),
			})
			return presenter.SendError(c, fiber.StatusNotFound, apperrors.ErrQuestionnaireNotFound.Error())
		}
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
	}

	if isAnonymous {
		// Check if the user has privileges to start the questionnaire
		hasPrivilege, err := h.roleService.HasPrivilegesOnInsance(ctx, c.UserContext(), req.UserID, req.QuestionnaireID, privilegeconstants.StartQuestionnariInsance)
		if err != nil {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return presenter.SendError(c, fiber.StatusNotFound, "questionnaire not found")
			}
			return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
		}
		if !hasPrivilege {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: logmessages.LogLackOfAuthorization,
			})
			return presenter.SendError(c, fiber.StatusForbidden, apperrors.ErrLackOfAuthorization.Error())
		}
	}

	// Check if the questionnaire is active
	isActive, err := h.questionnaireService.IsQuestionnaireActive(ctx, c.UserContext(), req.QuestionnaireID)
	if err != nil {
		if errors.Is(err, apperrors.ErrQuestionnaireNotFound) {
			logger.GetLogger().LogWarningFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: apperrors.ErrQuestionsNotFound.Error(),
			})
			return presenter.SendError(c, fiber.StatusNotFound, apperrors.ErrQuestionnaireNotFound.Error())
		}
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
	}
	if !isActive {
		return presenter.SendError(c, fiber.StatusBadRequest, "questionnaire is not active at this time")
	}

	// Call core service to start questionnaire
	submissionID, question, err := h.coreService.Start(ctx, c.UserContext(), req.UserID, req.QuestionnaireID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c, fiber.StatusNotFound, apperrors.ErrQuestionnaireNotFound.Error())
		}
		if errors.Is(err, apperrors.ErrQuestionsNotFound) {
			return presenter.SendError(c, fiber.StatusNotFound, apperrors.ErrQuestionsNotFound.Error())
		}
		if errors.Is(err, apperrors.ErrSubmissionLimit) {
			return presenter.SendError(c, fiber.StatusForbidden, apperrors.ErrSubmissionLimit.Error())
		}
		return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
	}

	// Prepare and send response
	resp := presenter.NewStartResponse(submissionID, question)
	return presenter.Send(c, fiber.StatusOK, true, "Questionnaire started", resp, nil)
}

// SubmitHandler receives an answer for the current question.
func (h *CoreHandler) SubmitHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: "submitting an answer",
	})

	// Extract user ID from context
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: "user_id not found in context",
		})
		return presenter.SendError(c, fiber.StatusUnauthorized, "Unauthorized: User ID missing")
	}

	// Parse and validate request
	req := &presenter.SubmitRequest{}
	if err := req.ParseAndValidate(c); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Set user ID in the request
	req.UserID = userID

	// Convert to domain and call service
	answerDomain := req.ToDomain()
	if err := h.coreService.Submit(ctx, c.UserContext(), req.SubmissionID, req.QuestionID, answerDomain); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c, fiber.StatusNotFound, apperrors.ErrNotFound.Error())
		}
		return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "Answer submitted successfully", nil, nil)
}

// BackHandler moves the current question index backwards (if allowed)
func (h *CoreHandler) BackHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse and validate request
	req := &presenter.NavigationRequest{}
	if err := req.ParseAndValidate(c); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Call service to move back
	question, err := h.coreService.Back(ctx, c.UserContext(), req.SubmissionID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		if errors.Is(err, apperrors.ErrBackIsNotAllowed) {
			presenter.SendError(c, fiber.StatusMethodNotAllowed, apperrors.ErrBackIsNotAllowed.Error())
		}
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Prepare and send response
	resp := presenter.NewGetQuestionResponse(question)
	return presenter.Send(c, fiber.StatusOK, true, "Moved back successfully", resp, nil)
}

// NextHandler moves the current question index forward to the next question.
func (h *CoreHandler) NextHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse and validate request
	req := &presenter.NavigationRequest{}
	if err := req.ParseAndValidate(c); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Call service to move next
	question, err := h.coreService.Next(ctx, c.UserContext(), req.SubmissionID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Prepare and send response
	resp := presenter.NewGetQuestionResponse(question)
	return presenter.Send(c, fiber.StatusOK, true, "Moved to next question", resp, nil)
}

// EndHandler finalizes the submission.
func (h *CoreHandler) EndHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse and validate request
	req := &presenter.EndRequest{}
	if err := req.ParseAndValidate(c); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Call service to end the questionnaire
	if err := h.coreService.End(ctx, c.UserContext(), req.SubmissionID); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "Questionnaire ended successfully", nil, nil)
}

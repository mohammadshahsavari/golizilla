package handler

import (
	"errors"
	"golizilla/handler/presenter"
	"golizilla/internal/apperrors"
	"golizilla/internal/logmessages"
	"golizilla/persistence/logger"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CoreHandler struct {
	coreService service.ICoreService
}

func NewCoreHandler(coreService service.ICoreService) *CoreHandler {
	return &CoreHandler{
		coreService: coreService,
	}
}

// StartHandler initializes a questionnaire session for the user and returns the first question.
func (h *CoreHandler) StartHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: "starting questionnaire",
	})

	questionnaireID, err := uuid.Parse(c.Params("questionnaire_id"))
	if err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid questionnaire_id format")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
	}

	submissionID, question, err := h.coreService.Start(ctx, c.UserContext(), userID, questionnaireID)
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

	// Now we return submission_id along with question details
	response := map[string]interface{}{
		"submission_id":    submissionID,
		"id":               question.ID,
		"questionnaire_id": question.QuestionnaireId,
		"index":            question.Index,
		"question_text":    question.QuestionText,
		"descriptive":      question.Descriptive,
	}

	return presenter.Send(c, fiber.StatusOK, true, "Questionnaire started", response, nil)
}

// SubmitHandler receives an answer for the current question.
func (h *CoreHandler) SubmitHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: "submitting an answer",
	})

	var req struct {
		SubmissionID string     `json:"submission_id"`
		QuestionID   string     `json:"question_id"`
		Descriptive  bool       `json:"descriptive"`
		Text         *string    `json:"text,omitempty"`
		OptionID     *uuid.UUID `json:"option_id,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	submissionID, err := uuid.Parse(req.SubmissionID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid submission_id format")
	}

	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid question_id format")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: "user id not found",
		})
		return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
	}

	createReq := presenter.CreateAnswerRequest{
		QuestionID:  questionID,
		Descriptive: req.Descriptive,
		Text:        req.Text,
		OptionID:    req.OptionID,
	}

	if err := createReq.Validate(); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	answerDomain := createReq.ToDomain(userID, submissionID)

	// Here we assume that the core service handles answer creation as part of submission
	// If your design differs (like using IAnswerService), adapt accordingly.
	err = h.coreService.Submit(ctx, c.UserContext(), submissionID, questionID, answerDomain)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c, fiber.StatusNotFound, apperrors.ErrNotFound.Error())
		}
		return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "answer submitted successfully", nil, nil)
}

// BackHandler moves the current question index backwards (if allowed)
func (h *CoreHandler) BackHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	var req struct {
		SubmissionID string `json:"submission_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	submissionID, err := uuid.Parse(req.SubmissionID)
	if err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid submission_id format")
	}

	question, err := h.coreService.Back(ctx, c.UserContext(), submissionID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})

		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "moved back successfully", presenter.NewGetQuestionResponse(question), nil)
}

// NextHandler moves the current question index forward to the next question
func (h *CoreHandler) NextHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	var req struct {
		SubmissionID string `json:"submission_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	submissionID, err := uuid.Parse(req.SubmissionID)
	if err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid submission_id format")
	}

	question, err := h.coreService.Next(ctx, c.UserContext(), submissionID)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})

		return presenter.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "moved to next question", presenter.NewGetQuestionResponse(question), nil)
}

// EndHandler finalizes the submission
func (h *CoreHandler) EndHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	var req struct {
		SubmissionID string `json:"submission_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	submissionID, err := uuid.Parse(req.SubmissionID)
	if err != nil {
		return presenter.SendError(c, fiber.StatusBadRequest, "invalid submission_id format")
	}

	if err := h.coreService.End(ctx, c.UserContext(), submissionID); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})

		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return presenter.Send(c, fiber.StatusOK, true, "questionnaire ended successfully", nil, nil)
}

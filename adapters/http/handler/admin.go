package handler

import (
	"golizilla/adapters/http/handler/presenter"
	"golizilla/adapters/persistence/logger"
	"golizilla/core/service"
	"golizilla/internal/apperrors"
	"golizilla/internal/logmessages"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService service.IAdminService
}

func NewAdminHandler(adminService service.IAdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.Context()
	userCtx := c.UserContext()

	// log: begin
	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogAdminHandler,
		Message: logmessages.LogAdminGetAllUsersBegin,
	})

	// Get query parameters for page and pageSize
	page, err := strconv.Atoi(c.Query("page", "1")) // Default to page 1 if not provided
	if err != nil || page < 1 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid page number")
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "5")) // Default to 5 items per page
	if err != nil || pageSize < 1 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid page size")
	}

	users, err := h.adminService.GetAllUsers(ctx, userCtx, page, pageSize)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c, 
			fiber.StatusInternalServerError, 
			apperrors.ErrInternalServerError.Error())
	}

	// TODO: end log

	return presenter.Send(c,
		fiber.StatusOK,
		true,
		"Users successfully fetched",
		users,
		nil,
	)
}


func (h *AdminHandler) GetAllQuestions(c *fiber.Ctx) error {
	ctx := c.Context()
	userCtx := c.UserContext()

	// TODO: log begin

	// Get query parameters for page and pageSize
	page, err := strconv.Atoi(c.Query("page", "1")) // Default to page 1 if not provided
	if err != nil || page < 1 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid page number")
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "5")) // Default to 5 items per page
	if err != nil || pageSize < 1 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid page size")
	}

	// Get paginated questions
	users, err := h.adminService.GetAllQuestions(ctx, userCtx, page, pageSize)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c, 
			fiber.StatusInternalServerError, 
			apperrors.ErrInternalServerError.Error())
	}

	// TODO: end log

	return presenter.Send(c,
		fiber.StatusOK,
		true,
		"Questions successfully fetched", // TODO: message
		users,
		nil,
	)
}

func (h *AdminHandler) GetAllQuestionnaires(c *fiber.Ctx) error {
	ctx := c.Context()
	userCtx := c.UserContext()

	// TODO: log begin

	// Get query parameters for page and pageSize
	page, err := strconv.Atoi(c.Query("page", "1")) // Default to page 1 if not provided
	if err != nil || page < 1 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid page number")
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "5")) // Default to 5 items per page
	if err != nil || pageSize < 1 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid page size")
	}

	// Get paginated questionnaires
	users, err := h.adminService.GetAllQuestionnaires(ctx, userCtx, page, pageSize)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c, 
			fiber.StatusInternalServerError, 
			apperrors.ErrInternalServerError.Error())
	}

	// TODO: end log

	return presenter.Send(c,
		fiber.StatusOK,
		true,
		"Questionnaires successfully fetched", // TODO: message
		users,
		nil,
	)
}

func (h *AdminHandler) GetAllRoles(c *fiber.Ctx) error {
	ctx := c.Context()
	userCtx := c.UserContext()

	// TODO: begin log

	// Get query parameters for page and pageSize
	page, err := strconv.Atoi(c.Query("page", "1")) // Default to page 1 if not provided
	if err != nil || page < 1 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid page number")
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "5")) // Default to 5 items per page
	if err != nil || pageSize < 1 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid page size")
	}

	// Get paginated roles
	users, err := h.adminService.GetAllRoles(ctx, userCtx, page, pageSize)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c, 
			fiber.StatusInternalServerError, 
			apperrors.ErrInternalServerError.Error())
	}

	// TODO: end log

	return presenter.Send(c,
		fiber.StatusOK,
		true,
		"Roles successfully fetched", // TODO: message
		users,
		nil,
	)
}

func (h *AdminHandler) GetAnswersByUserIDAndQuestionnaireID(c *fiber.Ctx) error {
	ctx := c.Context()
	userCtx := c.UserContext()

	// TODO: begin log

	// Get parameters for userID and questionnaireID
	userID, err := uuid.Parse(c.Params("userID"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid userID")
	}

	questionnaireID, err := uuid.Parse(c.Params("questionnaireID"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid userID")
	}

	// Get query parameters for page and pageSize
	page, err := strconv.Atoi(c.Query("page", "1")) // Default to page 1 if not provided
	if err != nil || page < 1 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid page number")
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "5")) // Default to 5 items per page
	if err != nil || pageSize < 1 {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"Invalid page size")
	}

	// Get paginated answers
	users, err := h.adminService.GetAnswersByUserIDAndQuestionnaireID(
		ctx, userCtx, userID, questionnaireID, page, pageSize)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogAdminHandler,
			Message: err.Error()})
		return presenter.SendError(c, 
			fiber.StatusInternalServerError, 
			apperrors.ErrInternalServerError.Error())
	}

	// TODO: end log

	return presenter.Send(c,
		fiber.StatusOK,
		true,
		"User answers successfully fetched", // TODO: message
		users,
		nil,
	)
}

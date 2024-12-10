package handler

import (
	"context"
	"fmt"
	"golizilla/adapters/http/handler/presenter"
	"golizilla/adapters/persistence/logger"
	"golizilla/core/service"
	"golizilla/internal/apperrors"
	"golizilla/internal/logmessages"
	privilegeconstants "golizilla/internal/privilege"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionnaireHandler struct {
	questionnaireService service.IQuestionnaireService
	roleService          service.IRoleService
}

func NewQuestionnaireHandler(
	questionnaireService service.IQuestionnaireService,
	roleService service.IRoleService) *QuestionnaireHandler {
	return &QuestionnaireHandler{
		questionnaireService: questionnaireService,
		roleService:          roleService,
	}
}

func (q *QuestionnaireHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireCreateBegin,
	})

	var request presenter.CreateQuestionnaireRequest
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
	id, err := q.questionnaireService.Create(ctx, c.UserContext(), userModel)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	err = presenter.Send(c, fiber.StatusOK, true, "Questionnaire created successfully", presenter.NewCreateQuestionnaireResponse(id), nil)
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

func (q *QuestionnaireHandler) Delete(c *fiber.Ctx) error {
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
	err = q.questionnaireService.Delete(ctx, c.UserContext(), id)
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

func (q *QuestionnaireHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireUpdateBegin,
	})

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	var request presenter.UpdateQuestionnaireRequest
	request.ID = id
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

	hasPrivilege, err := q.roleService.HasPrivilegesOnInsance(ctx, c.UserContext(), c.Locals("user_id").(uuid.UUID), request.ID, privilegeconstants.UpdateQuestionnaireInstance)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrInternalServerError.Error())
	}
	if !hasPrivilege {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: logmessages.LogLackOfAuthorization,
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrLackOfAuthorization.Error())
	}

	// Map fields to update
	updateFields := request.ToDomain()

	if err := q.questionnaireService.Update(ctx, c.UserContext(), request.ID, updateFields); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	err = presenter.Send(c, fiber.StatusOK, true, "Updated", nil, nil)
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

func (q *QuestionnaireHandler) GetById(c *fiber.Ctx) error {
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
	questionnaire, err := q.questionnaireService.GetById(ctx, c.UserContext(), id)
	if err != nil {
		if err == apperrors.ErrQuestionnaireNotFound {
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
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: logmessages.LogCastUserIdError,
		})
		return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
	}
	if questionnaire.Anonymous {
		hasPrivilege, err := q.roleService.HasPrivilegesOnInsance(ctx, c.UserContext(), userID, questionnaire.Id, privilegeconstants.ViewQuestionnaireInstances)
		if err != nil {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			return presenter.SendError(c, fiber.StatusInternalServerError, err.Error())
		}
		if !hasPrivilege {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: logmessages.LogLackOfAuthorization,
			})
			return presenter.SendError(c, fiber.StatusInternalServerError, apperrors.ErrLackOfAuthorization.Error())
		}
	}
	err = presenter.Send(c, fiber.StatusOK, true, "", presenter.NewGetQuestionnaireResponse(questionnaire), nil)
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

	return presenter.Send(c, fiber.StatusOK, true, "", presenter.NewGetQuestionnaireResponse(questionnaire), nil)
}

func (q *QuestionnaireHandler) GetByOwnerId(c *fiber.Ctx) error {
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
	questionnaires, err := q.questionnaireService.GetByOwnerId(ctx, c.UserContext(), id)
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

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireGetByOwnerIdSuccessful,
	})

	return presenter.Send(c, fiber.StatusOK, true, "", presenter.NewGetQuestionnairesResponse(questionnaires), nil)
}

func (q *QuestionnaireHandler) GiveAcess(c *fiber.Ctx) error {
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
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: logmessages.LogCastUserIdError,
		})
		return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
	}
	isOwner, err := q.questionnaireService.IsOwner(ctx, c.UserContext(), id, userID)
	if err != nil {
		if err == apperrors.ErrQuestionsNotFound {
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

	if !isOwner {
		return presenter.SendError(c, fiber.StatusUnauthorized, "you arent owner of this questionnari")
	}
	var request presenter.GiveAcessRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	for _, user := range request.UserIDs {
		role, err := q.roleService.GetRoleByUserId(ctx, c.UserContext(), user)
		if err != nil {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			return presenter.SendError(c, fiber.StatusBadRequest, "failed to get role")
		}
		err = q.roleService.AddPrivilegeOnInstance(ctx, c.UserContext(), role.ID, id, request.Privileges...)
		if err != nil {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			return presenter.SendError(c, fiber.StatusBadRequest, "failed to add privilges on instance")
		}
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireGiveAccessSuccessful,
	})

	return presenter.Send(c, fiber.StatusOK, true, "", nil, nil)
}

func (q *QuestionnaireHandler) DeleteAcess(c *fiber.Ctx) error {
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
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: logmessages.LogCastUserIdError,
		})
		return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
	}
	isOwner, err := q.questionnaireService.IsOwner(ctx, c.UserContext(), id, userID)
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

	if !isOwner {
		return presenter.SendError(c, fiber.StatusInternalServerError, "you arent owner of this questionnari")
	}
	var request presenter.GiveAcessRequest
	if err := c.BodyParser(&request); err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		return presenter.SendError(c, fiber.StatusBadRequest, apperrors.ErrInvalidInput.Error())
	}

	for _, user := range request.UserIDs {
		role, err := q.roleService.GetRoleByUserId(ctx, c.UserContext(), user)
		if err != nil {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			return presenter.SendError(c, fiber.StatusBadRequest, "failed to get role")
		}
		err = q.roleService.DeletePrivilegeOnInstance(ctx, c.UserContext(), role.ID, id, request.Privileges...)
		if err != nil {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			return presenter.SendError(c, fiber.StatusBadRequest, "failed to add privilges on instance")
		}
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireGiveAccessSuccessful,
	})

	return presenter.Send(c, fiber.StatusOK, true, "", nil, nil)
}

func (q *QuestionnaireHandler) GetResults(c *websocket.Conn) {
	ctx := context.Background()
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: logmessages.LogCastUserIdError,
		})
		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", apperrors.ErrInvalidUserID.Error())))
		return
	}

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
	isOwner, err := q.questionnaireService.IsOwner(ctx, nil, userID, id)
	if err != nil {
		logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
			Service: logmessages.LogQuestionnaireHandler,
			Message: err.Error(),
		})
		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
		return
	}
	if !isOwner {
		hasPrivilege, err := q.roleService.HasPrivilegesOnInsance(ctx, nil, userID, id, privilegeconstants.SeeResultsOnInstance)
		if err != nil {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
			return
		}
		if !hasPrivilege {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: logmessages.LogLackOfAuthorization,
			})
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", apperrors.ErrLackOfAuthorization.Error())))
			return
		}
	}
	_, err = q.questionnaireService.GetById(context.Background(), nil, id)
	if err != nil {
		if err == apperrors.ErrQuestionsNotFound {
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
		questionnaire, err := q.questionnaireService.GetById(context.Background(), nil, id)
		if err != nil {
			logger.GetLogger().LogErrorFromContext(ctx, logger.LogFields{
				Service: logmessages.LogQuestionnaireHandler,
				Message: err.Error(),
			})
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", err)))
			break
		}
		if lastValue != questionnaire.ParticipationCount {
			lastValue = questionnaire.ParticipationCount
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", lastValue)))
		}
		time.Sleep(time.Second * 10)
	}

	logger.GetLogger().LogInfoFromContext(ctx, logger.LogFields{
		Service: logmessages.LogQuestionnaireHandler,
		Message: logmessages.LogQuestionnaireGetResultsEnd,
	})
}

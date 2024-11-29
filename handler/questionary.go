package handler

import (
	"errors"
	"golizilla/handler/presenter"
	"golizilla/service"

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
		return presenter.BadRequest(c, err)
	}

	if err := request.Validate(); err != nil {
		return presenter.BadRequest(c, err)
	}

	userModel := request.ToDomain()

	if id, err := q.questionnaryService.Create(userModel); err != nil {
		//log
		return presenter.InternalServerError(c, err)
	} else {
		return presenter.Created(c, "Questionnary created successfully", presenter.NewCreateQuestionnaryResponse(id))
	}
}

func (q *QuestionnaryHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.BadRequest(c, errors.New("invalid ID format"))
	}
	err = q.questionnaryService.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presenter.NotFound(c, err)
		}

		//log
		return presenter.InternalServerError(c, err)
	}

	return presenter.OK(c, "Deleted", nil)
}

func (q *QuestionnaryHandler) Update(c *fiber.Ctx) error {
	var request presenter.CreateQuestionnaryRequest
	if err := c.BodyParser(&request); err != nil {
		return presenter.BadRequest(c, err)
	}

	if err := request.Validate(); err != nil {
		return presenter.BadRequest(c, err)
	}

	userModel := request.ToDomain()

	if err := q.questionnaryService.Update(userModel); err != nil {
		//log
		return presenter.InternalServerError(c, err)
	}

	return presenter.OK(c, "Updated", nil)
}

func (q *QuestionnaryHandler) GetById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.BadRequest(c, errors.New("invalid ID format"))
	}
	questionary, err := q.questionnaryService.GetById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presenter.NotFound(c, err)
		}

		//log
		return presenter.InternalServerError(c, err)
	}

	return presenter.OK(c, "", presenter.NewGetQuestionnaryResponse(questionary))
}

func (q *QuestionnaryHandler) GetByOwnerId(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return presenter.BadRequest(c, errors.New("invalid ID format"))
	}
	questionaries, err := q.questionnaryService.GetByOwnerId(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presenter.NotFound(c, err)
		}

		//log
		return presenter.InternalServerError(c, err)
	}

	return presenter.OK(c, "", presenter.NewGetQuestionnariesResponse(questionaries))
}

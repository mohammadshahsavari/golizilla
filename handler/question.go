package handler

import (
	"errors"
	"fmt"
	"golizilla/handler/presenter"
	"golizilla/internal/apperrors"
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

	var request presenter.CreateQuestionRequest
	if err := c.BodyParser(&request); err != nil {
		// log
		fmt.Printf("[Create Question] bad request error: %v", err)
		return presenter.SendError(c, 
			fiber.StatusBadRequest, 
			apperrors.ErrInvalidInput.Error(),
		)
	}

	if err := request.Validate(); err != nil {
		// log
		fmt.Printf("[Create Question] invalid input error: %v", err)
		return presenter.SendError(c, 
			fiber.StatusBadRequest, 
			apperrors.ErrInvalidInput.Error(),
		)
	}

	question := request.ToDomain()

	id, err := h.QuestionService.Create(ctx, question)
	if err != nil {
		// log
		fmt.Printf("[Create Question] Internal error: %v", err)
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

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		// log
		fmt.Printf("[Update Question] invalid input error: %v", err)
		return presenter.SendError(c, 
			fiber.StatusBadRequest, 
			"invalid ID format",
		)
	}

	var request presenter.UpdateQuestionRequest
	if err := c.BodyParser(&request); err != nil {
		// log
		fmt.Printf("[Update Question] invalid input error: %v", err)
		return presenter.SendError(c, 
			fiber.StatusBadRequest, 
			apperrors.ErrInvalidInput.Error(),
		)
	}

	if err := request.Validate(); err != nil {
		// log
		fmt.Printf("[Update Question] invalid input error: %v", err)
		return presenter.SendError(c, 
			fiber.StatusBadRequest, 
			apperrors.ErrInvalidInput.Error(),
		)
	}

	question := request.ToDomain()
	question.ID = id

	err = h.QuestionService.Update(ctx, question); 
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c, 
				fiber.StatusNotFound, 
				apperrors.ErrNotFound.Error(),
			)
		}
		// log
		fmt.Printf("[Create Question] internal error: %v", err)
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

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		// log
		return presenter.SendError(c, 
			fiber.StatusBadRequest, 
			"invalid ID format",
		)
	}

	err = h.QuestionService.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c, 
				fiber.StatusNotFound, 
				apperrors.ErrNotFound.Error(),
			)
		}
		// log
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

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		// log
		return presenter.SendError(c, 
			fiber.StatusBadRequest, 
			"invalid ID format",
		)
	}

	question, err := h.QuestionService.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c, 
				fiber.StatusNotFound, 
				apperrors.ErrNotFound.Error(),
			)
		}
		//log
		fmt.Printf("[Get Question] internal error: %v", err)
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

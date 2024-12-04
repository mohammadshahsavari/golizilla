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

	var request presenter.CreateAnswerRequest
	if err := c.BodyParser(&request); err != nil {
		// log
		fmt.Printf("[Create Answer] bad request error: %v", err)
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	if err := request.Validate(); err != nil {
		// log
		fmt.Printf("[Create Answer] invalid input error: %v", err)
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	Answer := request.ToDomain()

	id, err := h.answerService.Create(ctx, Answer)
	if err != nil {
		// log
		fmt.Printf("[Create Answer] Internal error: %v", err)
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

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		// log
		fmt.Printf("[Update Answer] invalid input error: %v", err)
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	var request presenter.UpdateAnswerRequest
	if err := c.BodyParser(&request); err != nil {
		// log
		fmt.Printf("[Update Answer] invalid input error: %v", err)
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	if err := request.Validate(); err != nil {
		// log
		fmt.Printf("[Update Answer] invalid input error: %v", err)
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			apperrors.ErrInvalidInput.Error(),
		)
	}

	Answer := request.ToDomain()
	Answer.ID = id

	err = h.answerService.Update(ctx, Answer)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c,
				fiber.StatusNotFound,
				apperrors.ErrNotFound.Error(),
			)
		}
		// log
		fmt.Printf("[Create Answer] internal error: %v", err)
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

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		// log
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	err = h.answerService.Delete(ctx, id)
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
		fmt.Sprintf("Answer with id: (%v) successfully deleted", id),
		nil,
		nil,
	)
}

func (h *AnswerHandler) GetByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		// log
		return presenter.SendError(c,
			fiber.StatusBadRequest,
			"invalid ID format",
		)
	}

	Answer, err := h.answerService.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenter.SendError(c,
				fiber.StatusNotFound,
				apperrors.ErrNotFound.Error(),
			)
		}
		//log
		fmt.Printf("[Get Answer] internal error: %v", err)
		return presenter.SendError(c,
			fiber.StatusInternalServerError,
			apperrors.ErrInternalServerError.Error(),
		)
	}

	return presenter.Send(c,
		fiber.StatusOK,
		true,
		"Answer successfully fetched",
		presenter.NewGetAnswerResponse(Answer),
		nil,
	)
}

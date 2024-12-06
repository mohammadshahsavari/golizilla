package handler

import (
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
)

type CoreHandler struct {
	coreService service.ICoreService
}

func NewCoreHandler(coreService service.ICoreService) *CoreHandler {
	return &CoreHandler{
		coreService: coreService,
	}
}

func (h *CoreHandler) StartHandler(c *fiber.Ctx) error {
	return nil
}
func (h *CoreHandler) SubmitHandler(c *fiber.Ctx) error {
	return nil
}
func (h *CoreHandler) BackHandler(c *fiber.Ctx) error {
	return nil
}
func (h *CoreHandler) NextHandler(c *fiber.Ctx) error {
	return nil
}
func (h *CoreHandler) EndHandler(c *fiber.Ctx) error {
	return nil
}

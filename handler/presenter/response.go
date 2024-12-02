package presenter

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// Response defines a standard API response format.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Send sends a custom response with a specified HTTP status code.
func Send(c *fiber.Ctx, statusCode int, success bool, message string, data interface{}, err error) error {
	response := Response{
		Success: success,
		Message: message,
		Data:    data,
	}

	if err != nil {
		response.Error = err.Error()
	}

	return c.Status(statusCode).JSON(response)
}

// SendError sends an error response with a specified HTTP status code.
func SendError(c *fiber.Ctx, statusCode int, errorMessage string) error {
	return Send(c, statusCode, false, "", nil, errors.New(errorMessage))
}

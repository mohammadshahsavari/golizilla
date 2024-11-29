package presenter

import (
	"github.com/gofiber/fiber/v2"
)

const (
	PostgresUniqueViolationCode = "23505"
)

// Response defines a standard API response format.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// NewResponse creates a base response with success set to true.
func NewResponse() *Response {
	return &Response{
		Success: true,
	}
}

// SetMessage sets the response message.
func (r *Response) SetMessage(message string) *Response {
	r.Message = message
	return r
}

// SetData sets the response data.
func (r *Response) SetData(data interface{}) *Response {
	r.Data = data
	return r
}

// SetError sets the response as a failure with an error message.
func (r *Response) SetError(err error) *Response {
	r.Success = false
	if err != nil {
		r.Error = err.Error()
	}
	return r
}

func SendError(c *fiber.Ctx, statusCode int, errorMessage string) error {
	return c.Status(statusCode).JSON(ErrorResponse{
		Success: false,
		Error:   errorMessage,
	})
}

// Send sends a custom response with a specified HTTP status code.
func Send(c *fiber.Ctx, status int, response *Response) error {
	return c.Status(status).JSON(response)
}

// OK sends a 200 OK response with a message and optional data.
func OK(c *fiber.Ctx, message string, data interface{}) error {
	return Send(c, fiber.StatusOK, NewResponse().SetMessage(message).SetData(data))
}

// Created sends a 201 Created response with a message and optional data.
func Created(c *fiber.Ctx, message string, data interface{}) error {
	return Send(c, fiber.StatusCreated, NewResponse().SetMessage(message).SetData(data))
}

// NoContent sends a 204 No Content response without any body.
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// BadRequest sends a 400 Bad Request response with an error message.
func BadRequest(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusBadRequest, NewResponse().SetError(err))
}

// Unauthorized sends a 401 Unauthorized response with an error message.
func Unauthorized(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusUnauthorized, NewResponse().SetError(err))
}

// Forbidden sends a 403 Forbidden response with an error message.
func Forbidden(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusForbidden, NewResponse().SetError(err))
}

// NotFound sends a 404 Not Found response with an error message.
func NotFound(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusNotFound, NewResponse().SetError(err))
}

// InternalServerError sends a 500 Internal Server Error response with an error message.
func InternalServerError(c *fiber.Ctx, err error) error {
	return Send(c, fiber.StatusInternalServerError, NewResponse().SetError(err))
}

// DuplicateEntry sends a 409 Conflict response for duplicate entry errors.
func DuplicateEntry(c *fiber.Ctx, err error) error {
	if err == nil {
		err = fiber.ErrConflict // Default error message for conflict
	}
	return Send(c, fiber.StatusConflict, NewResponse().SetError(err).SetMessage("Duplicate entry: the resource already exists"))
}

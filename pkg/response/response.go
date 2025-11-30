package response

import "github.com/gofiber/fiber/v2"

// Response is the standard structure for API responses
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success returns a successful response
func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created returns a response when creation is successful
func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error returns an error response
func Error(c *fiber.Ctx, statusCode int, message string, err error) error {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	return c.Status(statusCode).JSON(response)
}

// BadRequest returns a 400 error
func BadRequest(c *fiber.Ctx, message string, err error) error {
	return Error(c, fiber.StatusBadRequest, message, err)
}

// NotFound returns a 404 error
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, nil)
}

// InternalServerError returns a 500 error
func InternalServerError(c *fiber.Ctx, message string, err error) error {
	return Error(c, fiber.StatusInternalServerError, message, err)
}

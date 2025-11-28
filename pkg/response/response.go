package response

import "github.com/gofiber/fiber/v2"

// Response là cấu trúc chuẩn cho API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success trả về response thành công
func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created trả về response khi tạo mới thành công
func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error trả về response lỗi
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

// BadRequest trả về lỗi 400
func BadRequest(c *fiber.Ctx, message string, err error) error {
	return Error(c, fiber.StatusBadRequest, message, err)
}

// NotFound trả về lỗi 404
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, nil)
}

// InternalServerError trả về lỗi 500
func InternalServerError(c *fiber.Ctx, message string, err error) error {
	return Error(c, fiber.StatusInternalServerError, message, err)
}

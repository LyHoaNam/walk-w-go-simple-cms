package middleware

import (
	"simple-template/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ErrorHandler middleware handles global errors
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Catch and handle panics
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Panic recovered: %v", r)
				_ = response.InternalServerError(c, "Internal server error", nil)
			}
		}()

		// Process request
		err := c.Next()

		// If there's an error, handle it here
		if err != nil {
			// Fiber error
			if e, ok := err.(*fiber.Error); ok {
				return response.Error(c, e.Code, e.Message, nil)
			}

			// Lỗi khác
			log.Errorf("Unhandled error: %v", err)
			return response.InternalServerError(c, "Internal server error", err)
		}

		return nil
	}
}

package middleware

import (
	"simple-template/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ErrorHandler middleware xử lý lỗi toàn cục
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Bắt panic và xử lý
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Panic recovered: %v", r)
				_ = response.InternalServerError(c, "Internal server error", nil)
			}
		}()

		// Xử lý request
		err := c.Next()

		// Nếu có lỗi, xử lý ở đây
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

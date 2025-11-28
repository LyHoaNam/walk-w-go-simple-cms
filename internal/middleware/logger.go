package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// Logger middleware ghi log cho mỗi request
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Lưu thời gian bắt đầu
		start := time.Now()

		// Xử lý request
		err := c.Next()

		// Tính thời gian xử lý
		duration := time.Since(start)

		// Log thông tin request
		log.Infof(
			"[%s] %s %s - Status: %d - Duration: %v",
			c.IP(),
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}

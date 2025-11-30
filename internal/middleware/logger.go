package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// Logger middleware logs each request
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Save start time
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate processing time
		duration := time.Since(start)

		// Log request information
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

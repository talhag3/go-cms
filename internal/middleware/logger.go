package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Logger returns a configured request logger
func Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	})
}

// RequestID adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := time.Now().Format("20060102150405") + "-" + c.IP()
		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)
		return c.Next()
	}
}

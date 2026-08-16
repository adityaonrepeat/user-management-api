package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

const requestIDKey = "requestID"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}

		c.Locals(requestIDKey, id)
		c.Set(RequestIDHeader, id)

		return c.Next()
	}
}

func RequestIDFrom(c *fiber.Ctx) string {
	id, _ := c.Locals(requestIDKey).(string)
	return id
}

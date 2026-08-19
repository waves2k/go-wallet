package middleware

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type RequestIDKeyType string

const RequestIDKey RequestIDKeyType = "request_id"

func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestID := c.Get(string(RequestIDKey))
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set(string(RequestIDKey), requestID)
		}

		return c.Next()
	}
}

func RequestIDFromContext(c context.Context) any {
	return c.Value(RequestIDKey)
}

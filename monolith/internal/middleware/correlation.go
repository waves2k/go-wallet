package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func CorrelationID() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !c.HasHeader("X-Correlation-ID") {
			corID := uuid.New().String()
			c.Response().Header.Add("X-Correlation-ID", corID)
		}

		return c.Next()
	}
}

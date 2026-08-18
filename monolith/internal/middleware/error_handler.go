package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	appErr "github.com/waves2k/go-wallet/monolith/internal/errors"
	"github.com/waves2k/go-wallet/monolith/internal/logger"
)

var (
	ErrInternalServerFailure = errors.New("internal server failure")
)

func ErrorHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := c.Next(); err != nil {
			if appErr, ok := err.(*appErr.AppError); ok {
				logger.Warn(c.Context(), "Client error occured",
					"code", appErr.Code,
					"message", appErr.Message,
					"status", appErr.StatusCode,
				)

				return writeError(c, appErr.StatusCode, appErr)
			}
			logger.Error(c.Context(), "Unhandled error", "error", err)
			return writeError(c, fiber.StatusInternalServerError, ErrInternalServerFailure)
		}
		return nil
	}
}

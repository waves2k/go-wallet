package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/waves2k/go-wallet/monolith/internal/auth"
	appErr "github.com/waves2k/go-wallet/monolith/internal/errors"
)

type UserIDKeyType string

const UserIDKey UserIDKeyType = "user_id"

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return appErr.ErrUnauthorized
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		if tokenStr == header || tokenStr == "" {
			return appErr.ErrUnauthorized
		}

		claims, err := auth.ValidateToken(tokenStr)

		if err != nil {
			fmt.Printf(err.Error())
			return appErr.ErrUnauthorized
		}

		ctx := context.WithValue(c.Context(), UserIDKey, claims.ID)
		c.SetContext(ctx)
		return c.Next()
	}
}

func UserIDFromContext(c context.Context) any {
	return c.Value(UserIDKey)
}

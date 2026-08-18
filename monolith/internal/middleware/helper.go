package middleware

import "github.com/gofiber/fiber/v3"

func writeJSON(c fiber.Ctx, code int, v any) error {
	return c.Status(code).JSON(fiber.Map{
		"success": true,
		"result":  v,
	})
}

func writeError(c fiber.Ctx, code int, e error) error {
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": e.Error(),
	})
}

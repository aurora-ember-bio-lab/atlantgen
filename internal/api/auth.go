package api

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// APIKeyMiddleware validates Bearer token authentication.
// Set API_KEY env var to enable. If empty, auth is disabled (open access).
func APIKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			return c.Next()
		}

		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing Authorization header",
				"hint":  "Send: Authorization: Bearer <API_KEY>",
			})
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid Authorization format",
				"hint":  "Use: Authorization: Bearer <API_KEY>",
			})
		}

		if token != apiKey {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "invalid API key",
			})
		}

		return c.Next()
	}
}

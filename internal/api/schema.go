package api

import (
	"github.com/gofiber/fiber/v2"

	"aura-amber-saas/internal/migration"
)

type schemaRequest struct {
	TargetDBURL string `json:"targetDbUrl"`
}

// RegisterSchema wires the Atlas Schema Orchestrator endpoints.
func RegisterSchema(app *fiber.App) {
	app.Post("/api/schema/diff", schemaDiff)
	app.Post("/api/schema/apply", schemaApply)
}

func schemaDiff(c *fiber.Ctx) error {
	var req schemaRequest
	if err := c.BodyParser(&req); err != nil || req.TargetDBURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "targetDbUrl is required"})
	}

	diff, err := migration.DiffSchema(c.Context(), req.TargetDBURL, "")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"diff": diff})
}

func schemaApply(c *fiber.Ctx) error {
	var req schemaRequest
	if err := c.BodyParser(&req); err != nil || req.TargetDBURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "targetDbUrl is required"})
	}

	if err := migration.ApplySchema(c.Context(), req.TargetDBURL, ""); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "applied"})
}

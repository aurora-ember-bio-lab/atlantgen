package api

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"os"

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
	app.Get("/api/schema/inspect", schemaInspect)
}

func schemaDiff(c *fiber.Ctx) error {
	var req schemaRequest
	if err := c.BodyParser(&req); err != nil || req.TargetDBURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "targetDbUrl is required"})
	}

	// Use atlas schema inspect to compare current state
	cmd := exec.CommandContext(context.Background(), "atlas", "schema", "inspect",
		"--url", req.TargetDBURL,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("atlas schema inspect failed: %v: %s", err, string(out))})
	}
	return c.JSON(fiber.Map{"schema": string(out)})
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

func schemaInspect(c *fiber.Ctx) error {
	targetDBURL := os.Getenv("DATABASE_URL")
	if targetDBURL == "" {
		targetDBURL = "postgres://aura:aura@localhost:5432/aura_amber?sslmode=disable"
	}

	cmd := exec.CommandContext(context.Background(), "atlas", "schema", "inspect",
		"--url", targetDBURL,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("atlas inspect failed: %v: %s", err, stderr.String())})
	}
	return c.JSON(fiber.Map{"schema": stdout.String()})
}

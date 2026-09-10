package api

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"aura-amber-saas/internal/migration"
)

type createMigrationRequest struct {
	Name                 string            `json:"name"`
	SourceType           string            `json:"sourceType"`
	SourceConnectionInfo map[string]string `json:"sourceConnectionInfo"`
	TargetDBURL          string            `json:"targetDbUrl"`
}

// RegisterMigrations wires the migration-job endpoints.
func RegisterMigrations(app *fiber.App) {
	app.Get("/api/migrations", listMigrations)
	app.Get("/api/migrations/:id", getMigration)
	app.Post("/api/migrations", createMigration)
	app.Delete("/api/migrations/:id", deleteMigration)
}

func listMigrations(c *fiber.Ctx) error {
	return c.JSON(migration.ListJobs())
}

func getMigration(c *fiber.Ctx) error {
	id := c.Params("id")
	job := migration.GetJob(id)
	if job == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "migration not found"})
	}
	return c.JSON(job)
}

func createMigration(c *fiber.Ctx) error {
	var req createMigrationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Name == "" || req.SourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name and sourceType are required"})
	}

	// Default targetDbUrl to DATABASE_URL env var if not provided
	if req.TargetDBURL == "" {
		req.TargetDBURL = os.Getenv("DATABASE_URL")
	}
	if req.TargetDBURL == "" {
		req.TargetDBURL = "postgres://aura:aura@localhost:5432/aura_amber?sslmode=disable"
	}

	job := migration.EnqueueJob(migration.Job{
		Name:        req.Name,
		SourceType:  migration.SourceType(req.SourceType),
		TargetDBURL: req.TargetDBURL,
	})

	workerJobID, err := migration.TriggerWorker(c.Context(), req.Name, job.SourceType, req.SourceConnectionInfo, req.TargetDBURL)
	if err != nil {
		migration.MarkFailed(job.ID)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"job":   job,
			"error": "created job but could not reach the migration worker: " + err.Error(),
		})
	}
	migration.SetWorkerJobID(job.ID, workerJobID)

	return c.Status(fiber.StatusCreated).JSON(job)
}

func deleteMigration(c *fiber.Ctx) error {
	id := c.Params("id")
	if migration.DeleteJob(id) {
		return c.JSON(fiber.Map{"status": "deleted", "id": id})
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "migration not found"})
}

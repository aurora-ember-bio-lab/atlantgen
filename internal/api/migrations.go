package api

import (
	"github.com/gofiber/fiber/v2"

	"aura-amber-saas/internal/migration"
)

type createMigrationRequest struct {
	Name                 string            `json:"name"`
	SourceType           string            `json:"sourceType"`
	SourceConnectionInfo map[string]string `json:"sourceConnectionInfo"`
	TargetDBURL          string            `json:"targetDbUrl"`
}

// RegisterMigrations wires the migration-job endpoints. Creating a job here
// also hands it off to the worker (worker/src/index.ts) over HTTP so it
// actually runs — see internal/migration/worker_client.go.
func RegisterMigrations(app *fiber.App) {
	app.Get("/api/migrations", listMigrations)
	app.Post("/api/migrations", createMigration)
}

func listMigrations(c *fiber.Ctx) error {
	return c.JSON(migration.ListJobs())
}

func createMigration(c *fiber.Ctx) error {
	var req createMigrationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Name == "" || req.SourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name and sourceType are required"})
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

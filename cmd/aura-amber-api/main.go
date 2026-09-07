package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"aura-amber-saas/internal/api"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api.RegisterMigrations(app)
	api.RegisterSchema(app)
	api.RegisterGitHubWebhooks(app)

	log.Fatal(app.Listen(":8080"))
}

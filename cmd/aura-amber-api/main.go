package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"

	"aura-amber-saas/internal/api"
	"aura-amber-saas/internal/migration"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize database connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		dbURL = "postgres://aura:aura@localhost:5432/aura_amber?sslmode=disable"
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("WARNING: failed to connect to database: %v (running in degraded mode)", err)
	} else {
		if err := conn.Ping(); err != nil {
			log.Printf("WARNING: database ping failed: %v (running in degraded mode)", err)
		} else {
			log.Println("Connected to PostgreSQL")
			migration.InitDB(conn)
			api.InitAccountsDB(conn)
		}
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api.RegisterMigrations(app)
	api.RegisterConnectors(app)
	api.RegisterSchema(app)
	api.RegisterGitHubWebhooks(app)
	api.RegisterBilling(app)

	log.Fatal(app.Listen(":8080"))
}

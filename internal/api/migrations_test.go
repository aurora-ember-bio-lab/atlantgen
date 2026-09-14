package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupMigrationsApp() *fiber.App {
	app := fiber.New()
	RegisterMigrations(app)
	return app
}

func TestCreateMigration_MissingFields(t *testing.T) {
	app := setupMigrationsApp()

	req := httptest.NewRequest(http.MethodPost, "/api/migrations", strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing fields, got %d", resp.StatusCode)
	}
}

func TestCreateMigration_InvalidBody(t *testing.T) {
	app := setupMigrationsApp()

	req := httptest.NewRequest(http.MethodPost, "/api/migrations", strings.NewReader(`{invalid json`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid body, got %d", resp.StatusCode)
	}
}

func TestCreateMigration_NoWorker(t *testing.T) {
	// Worker is not reachable in tests → expect 502 and a created job
	os.Setenv("WORKER_URL", "http://localhost:1") // unreachable
	defer os.Unsetenv("WORKER_URL")
	os.Setenv("DATABASE_URL", "postgres://x:x@localhost:1/db")
	defer os.Unsetenv("DATABASE_URL")

	app := setupMigrationsApp()
	req := httptest.NewRequest(http.MethodPost, "/api/migrations", strings.NewReader(`{"name":"test-mig","sourceType":"woocommerce"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected 502 when worker unreachable, got %d", resp.StatusCode)
	}
}

func TestGetMigration_NotFound(t *testing.T) {
	app := setupMigrationsApp()

	req := httptest.NewRequest(http.MethodGet, "/api/migrations/nonexistent", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown migration, got %d", resp.StatusCode)
	}
}

func TestDeleteMigration_NotFound(t *testing.T) {
	app := setupMigrationsApp()

	req := httptest.NewRequest(http.MethodDelete, "/api/migrations/nonexistent", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown migration delete, got %d", resp.StatusCode)
	}
}

func TestListMigrations_Empty(t *testing.T) {
	app := setupMigrationsApp()

	req := httptest.NewRequest(http.MethodGet, "/api/migrations", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for empty list, got %d", resp.StatusCode)
	}
}
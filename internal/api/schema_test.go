package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupSchemaApp() *fiber.App {
	app := fiber.New()
	RegisterSchema(app)
	return app
}

func TestSchemaDiff_MissingURL(t *testing.T) {
	app := setupSchemaApp()

	req := httptest.NewRequest(http.MethodPost, "/api/schema/diff", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing targetDbUrl, got %d", resp.StatusCode)
	}
}

func TestSchemaDiff_InvalidBody(t *testing.T) {
	app := setupSchemaApp()

	req := httptest.NewRequest(http.MethodPost, "/api/schema/diff", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid body, got %d", resp.StatusCode)
	}
}

func TestSchemaApply_MissingURL(t *testing.T) {
	app := setupSchemaApp()

	req := httptest.NewRequest(http.MethodPost, "/api/schema/apply", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing targetDbUrl, got %d", resp.StatusCode)
	}
}
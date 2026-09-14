package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupAuthApp() *fiber.App {
	app := fiber.New()
	app.Use(APIKeyMiddleware())
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})
	return app
}

func TestAPIKeyMiddleware_DisabledWhenKeyEmpty(t *testing.T) {
	os.Setenv("API_KEY", "")
	app := setupAuthApp()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when auth disabled, got %d", resp.StatusCode)
	}
}

func TestAPIKeyMiddleware_MissingHeader(t *testing.T) {
	os.Setenv("API_KEY", "secret123")
	defer os.Unsetenv("API_KEY")
	app := setupAuthApp()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing header, got %d", resp.StatusCode)
	}
}

func TestAPIKeyMiddleware_WrongKey(t *testing.T) {
	os.Setenv("API_KEY", "secret123")
	defer os.Unsetenv("API_KEY")
	app := setupAuthApp()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer wrongkey")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for wrong key, got %d", resp.StatusCode)
	}
}

func TestAPIKeyMiddleware_MalformedHeader(t *testing.T) {
	os.Setenv("API_KEY", "secret123")
	defer os.Unsetenv("API_KEY")
	app := setupAuthApp()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Basic foo") // not Bearer
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for malformed header, got %d", resp.StatusCode)
	}
}

func TestAPIKeyMiddleware_CorrectKey(t *testing.T) {
	os.Setenv("API_KEY", "secret123")
	defer os.Unsetenv("API_KEY")
	app := setupAuthApp()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer secret123")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for correct key, got %d", resp.StatusCode)
	}
}
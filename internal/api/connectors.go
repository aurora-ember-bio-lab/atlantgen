package api

import (
	"github.com/gofiber/fiber/v2"
)

type ConnectorState struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	LastSync  string `json:"lastSync,omitempty"`
	Connected bool   `json:"connected"`
}

var connectors = []ConnectorState{
	{ID: "wordpress", Name: "WordPress", State: "connected", LastSync: "2026-09-05T14:12:00Z", Connected: true},
	{ID: "woocommerce", Name: "WooCommerce", State: "connected", LastSync: "2026-09-05T14:12:00Z", Connected: true},
	{ID: "shopify", Name: "Shopify", State: "not_connected", Connected: false},
	{ID: "magento", Name: "Magento", State: "error", LastSync: "2026-09-02T18:05:00Z", Connected: false},
}

func RegisterConnectors(app *fiber.App) {
	app.Get("/api/connectors", listConnectors)
	app.Put("/api/connectors/:id", updateConnector)
}

func listConnectors(c *fiber.Ctx) error {
	return c.JSON(connectors)
}

func updateConnector(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		State string `json:"state"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	for i, cn := range connectors {
		if cn.ID == id {
			connectors[i].State = req.State
			if req.State == "connected" {
				connectors[i].Connected = true
				connectors[i].LastSync = "now"
			} else {
				connectors[i].Connected = false
			}
			return c.JSON(connectors[i])
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "connector not found"})
}

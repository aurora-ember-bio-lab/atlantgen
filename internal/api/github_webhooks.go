package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
)

// RegisterGitHubWebhooks wires the GitHub App webhook endpoint: Marketplace
// purchase events (billing, no Stripe involved) and installation lifecycle
// events. Configure the webhook URL and secret per github-app/manifest.yml.
func RegisterGitHubWebhooks(app *fiber.App) {
	app.Post("/github/webhooks", handleGitHubWebhook)
}

type marketplacePurchaseEvent struct {
	Action string `json:"action"` // purchased | cancelled | changed | pending_change | pending_change_cancelled
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
	MarketplacePurchase struct {
		Account struct {
			Login string `json:"login"`
			Type  string `json:"type"`
		} `json:"account"`
		BillingCycle string `json:"billing_cycle"`
		Plan         struct {
			Name                string `json:"name"`
			MonthlyPriceInCents int    `json:"monthly_price_in_cents"`
		} `json:"plan"`
	} `json:"marketplace_purchase"`
}

type installationEvent struct {
	Action       string `json:"action"` // created | deleted | suspend | unsuspend
	Installation struct {
		ID      int64 `json:"id"`
		Account struct {
			Login string `json:"login"`
		} `json:"account"`
	} `json:"installation"`
}

func handleGitHubWebhook(c *fiber.Ctx) error {
	body := c.Body()

	if secret := os.Getenv("GITHUB_WEBHOOK_SECRET"); secret != "" {
		if !verifySignature(secret, body, c.Get("X-Hub-Signature-256")) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid signature"})
		}
	}

	switch c.Get("X-GitHub-Event") {
	case "marketplace_purchase":
		return handleMarketplacePurchase(c, body)
	case "installation":
		return handleInstallation(c, body)
	default:
		// installation_repositories, ping, etc. — acknowledge, nothing to do yet.
		return c.SendStatus(fiber.StatusOK)
	}
}

func verifySignature(secret string, body []byte, header string) bool {
	if len(header) < 8 || header[:7] != "sha256=" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(header))
}

func handleMarketplacePurchase(c *fiber.Ctx, body []byte) error {
	var evt marketplacePurchaseEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	switch evt.Action {
	case "purchased":
		// TODO: activate evt.MarketplacePurchase.Account.Login on
		// evt.MarketplacePurchase.Plan.Name (persist to Postgres once the
		// accounts table exists).
	case "changed":
		// TODO: update the stored plan for the account.
	case "cancelled":
		// TODO: deactivate the account's access.
	case "pending_change", "pending_change_cancelled":
		// No-op — the change isn't in effect yet.
	}

	return c.SendStatus(fiber.StatusOK)
}

func handleInstallation(c *fiber.Ctx, body []byte) error {
	var evt installationEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	switch evt.Action {
	case "created":
		// TODO: record evt.Installation.ID for evt.Installation.Account.Login.
	case "deleted", "suspend":
		// TODO: revoke access for the installation's account.
	case "unsuspend":
		// TODO: restore access.
	}

	return c.SendStatus(fiber.StatusOK)
}

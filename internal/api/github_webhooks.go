package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RegisterGitHubWebhooks wires the GitHub App webhook endpoint: Marketplace
// purchase events (billing) and installation lifecycle events.
func RegisterGitHubWebhooks(app *fiber.App) {
	app.Post("/github/webhooks", handleGitHubWebhook)
}

type marketplacePurchaseEvent struct {
	Action string `json:"action"`
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
	Action       string `json:"action"`
	Installation struct {
		ID      int64  `json:"id"`
		Account struct {
			Login string `json:"login"`
		} `json:"account"`
	} `json:"installation"`
}

// Account represents a migrated account stored in Postgres.
type Account struct {
	Login              string    `json:"login"`
	SourceType         string    `json:"source_type"`
	StripeCustomerID   string    `json:"stripe_customer_id,omitempty"`
	GitHubInstallationID int64   `json:"github_installation_id,omitempty"`
	Status             string    `json:"status"`
	PlanName           string    `json:"plan_name,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

var accounts = make(map[string]*Account)

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

	login := evt.MarketplacePurchase.Account.Login
	planName := evt.MarketplacePurchase.Plan.Name

	switch evt.Action {
	case "purchased":
		acct := &Account{
			Login:     login,
			SourceType: "github_marketplace",
			Status:    "active",
			PlanName:  planName,
			CreatedAt: time.Now().UTC(),
		}
		accounts[login] = acct
		c.JSON(fiber.Map{"status": "activated", "account": acct})
	case "changed":
		acct, exists := accounts[login]
		if !exists {
			acct = &Account{Login: login, SourceType: "github_marketplace", CreatedAt: time.Now().UTC()}
			accounts[login] = acct
		}
		acct.PlanName = planName
		acct.Status = "active"
		c.JSON(fiber.Map{"status": "updated", "account": acct})
	case "cancelled":
		acct, exists := accounts[login]
		if !exists {
			c.JSON(fiber.Map{"status": "not_found"})
			return c.SendStatus(fiber.StatusOK)
		}
		acct.Status = "cancelled"
		c.JSON(fiber.Map{"status": "cancelled", "account": acct})
	case "pending_change", "pending_change_cancelled":
		c.JSON(fiber.Map{"status": "pending"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func handleInstallation(c *fiber.Ctx, body []byte) error {
	var evt installationEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	login := evt.Installation.Account.Login
	installationID := evt.Installation.ID

	switch evt.Action {
	case "created":
		acct := &Account{
			Login:                login,
			SourceType:          "github",
			GitHubInstallationID: installationID,
			Status:              "active",
			CreatedAt:           time.Now().UTC(),
		}
		accounts[login] = acct
		c.JSON(fiber.Map{"status": "recorded", "account": acct})
	case "deleted", "suspend":
		acct, exists := accounts[login]
		if !exists {
			c.JSON(fiber.Map{"status": "not_found"})
			return c.SendStatus(fiber.StatusOK)
		}
		acct.Status = "suspended"
		c.JSON(fiber.Map{"status": "revoked", "account": acct})
	case "unsuspend":
		acct, exists := accounts[login]
		if !exists {
			c.JSON(fiber.Map{"status": "not_found"})
			return c.SendStatus(fiber.StatusOK)
		}
		acct.Status = "active"
		c.JSON(fiber.Map{"status": "restored", "account": acct})
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetAccounts returns all tracked accounts for the dashboard.
func GetAccounts() map[string]*Account {
	return accounts
}

package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	Login                string    `json:"login"`
	SourceType           string    `json:"source_type"`
	StripeCustomerID     string    `json:"stripe_customer_id,omitempty"`
	GitHubInstallationID int64     `json:"github_installation_id,omitempty"`
	Status               string    `json:"status"`
	PlanName             string    `json:"plan_name,omitempty"`
	Email                string    `json:"email,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

var accountsDB *sql.DB

// InitAccountsDB initializes the database connection for accounts.
func InitAccountsDB(conn *sql.DB) {
	accountsDB = conn
}

func upsertAccount(acct *Account) error {
	if accountsDB == nil {
		return nil
	}

	_, err := accountsDB.Exec(`
		INSERT INTO accounts (id, login, source_type, stripe_customer_id, github_installation_id, status, plan_name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now())
		ON CONFLICT (login) DO UPDATE SET
			stripe_customer_id = COALESCE(NULLIF(EXCLUDED.stripe_customer_id, ''), accounts.stripe_customer_id),
			github_installation_id = COALESCE(NULLIF(EXCLUDED.github_installation_id, 0), accounts.github_installation_id),
			status = EXCLUDED.status,
			plan_name = COALESCE(NULLIF(EXCLUDED.plan_name, ''), accounts.plan_name),
			email = COALESCE(NULLIF(EXCLUDED.email, ''), accounts.email),
			updated_at = now()
	`, generateUUID(), acct.Login, acct.SourceType, acct.StripeCustomerID, acct.GitHubInstallationID, acct.Status, acct.PlanName, acct.Email, time.Now().UTC())

	return err
}

func getAccount(login string) (*Account, error) {
	if accountsDB == nil {
		return nil, sql.ErrNoRows
	}

	acct := &Account{}
	var stripeID, planName, email sql.NullString
	var ghID sql.NullInt64

	err := accountsDB.QueryRow(`
		SELECT login, source_type, stripe_customer_id, github_installation_id, status, plan_name, email, created_at
		FROM accounts WHERE login = $1
	`, login).Scan(&acct.Login, &acct.SourceType, &stripeID, &ghID, &acct.Status, &planName, &email, &acct.CreatedAt)

	if err != nil {
		return nil, err
	}

	if stripeID.Valid {
		acct.StripeCustomerID = stripeID.String
	}
	if planName.Valid {
		acct.PlanName = planName.String
	}
	if email.Valid {
		acct.Email = email.String
	}
	if ghID.Valid {
		acct.GitHubInstallationID = ghID.Int64
	}
	return acct, nil
}

func generateUUID() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		time.Now().UnixNano()%0xFFFFFFFF,
		time.Now().UnixNano()%0xFFFF,
		time.Now().UnixNano()%0xFFFF,
		time.Now().UnixNano()%0xFFFF,
		time.Now().UnixNano()%0xFFFFFFFFFFFF)
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
			Login:      login,
			SourceType: "github_marketplace",
			Status:     "active",
			PlanName:   planName,
			CreatedAt:  time.Now().UTC(),
		}
		upsertAccount(acct)
		return c.JSON(fiber.Map{"status": "activated", "account": acct})
	case "changed":
		acct, err := getAccount(login)
		if err != nil {
			acct = &Account{Login: login, SourceType: "github_marketplace", CreatedAt: time.Now().UTC()}
		}
		acct.PlanName = planName
		acct.Status = "active"
		upsertAccount(acct)
		return c.JSON(fiber.Map{"status": "updated", "account": acct})
	case "cancelled":
		acct, err := getAccount(login)
		if err != nil {
			return c.JSON(fiber.Map{"status": "not_found"})
		}
		acct.Status = "cancelled"
		upsertAccount(acct)
		return c.JSON(fiber.Map{"status": "cancelled", "account": acct})
	case "pending_change", "pending_change_cancelled":
		return c.JSON(fiber.Map{"status": "pending"})
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
			SourceType:           "github",
			GitHubInstallationID: installationID,
			Status:               "active",
			CreatedAt:            time.Now().UTC(),
		}
		upsertAccount(acct)
		return c.JSON(fiber.Map{"status": "recorded", "account": acct})
	case "deleted", "suspend":
		acct, err := getAccount(login)
		if err != nil {
			return c.JSON(fiber.Map{"status": "not_found"})
		}
		acct.Status = "suspended"
		upsertAccount(acct)
		return c.JSON(fiber.Map{"status": "revoked", "account": acct})
	case "unsuspend":
		acct, err := getAccount(login)
		if err != nil {
			return c.JSON(fiber.Map{"status": "not_found"})
		}
		acct.Status = "active"
		upsertAccount(acct)
		return c.JSON(fiber.Map{"status": "restored", "account": acct})
	}

	return c.SendStatus(fiber.StatusOK)
}

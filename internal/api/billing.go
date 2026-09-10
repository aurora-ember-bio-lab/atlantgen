package api

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

// RegisterBilling wires the Stripe billing endpoints.
func RegisterBilling(app *fiber.App) {
	app.Post("/api/billing/checkout", createCheckout)
	app.Post("/api/billing/customer", createCustomer)
	app.Get("/api/billing/account/:login", getBillingAccount)
}

// createCheckout creates a Stripe Checkout Session for a subscription.
// POST body: { login, plan, email }
func createCheckout(c *fiber.Ctx) error {
	var req struct {
		Login string `json:"login"`
		Plan  string `json:"plan"`
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.Login == "" || req.Plan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "login and plan are required"})
	}

	acct, exists := accounts[req.Login]
	if !exists {
		acct = &Account{
			Login:     req.Login,
			SourceType: "github_marketplace",
			Status:    "active",
			CreatedAt: time.Now().UTC(),
		}
		accounts[req.Login] = acct
	}

	priceID := os.Getenv("STRIPE_PRICE_ID_PLAN_BASIC")
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: []*string{stripe.String("card")},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String("subscription"),
		SuccessURL: stripe.String(os.Getenv("NEXT_PUBLIC_API_URL") + "/billing/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(os.Getenv("NEXT_PUBLIC_API_URL") + "/billing/cancel"),
	}

	if acct.StripeCustomerID != "" {
		params.Customer = stripe.String(acct.StripeCustomerID)
	}

	stripeSession, err := session.New(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":    "checkout_created",
		"sessionId": stripeSession.ID,
		"login":     req.Login,
		"plan":      req.Plan,
	})
}

// createCustomer creates a Stripe customer for a marketplace account.
// POST body: { login, email, plan }
func createCustomer(c *fiber.Ctx) error {
	var req struct {
		Login string `json:"login"`
		Email string `json:"email"`
		Plan  string `json:"plan"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.Login == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "login and email are required"})
	}

	acct, exists := accounts[req.Login]
	if !exists {
		acct = &Account{
			Login:     req.Login,
			SourceType: "github_marketplace",
			CreatedAt: time.Now().UTC(),
		}
		accounts[req.Login] = acct
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(req.Email),
		Name:  stripe.String(req.Login),
		Metadata: map[string]string{
			"login":  req.Login,
			"plan":   req.Plan,
			"source": "aura-amber",
		},
	}

	stripeCust, err := customer.New(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	acct.StripeCustomerID = stripeCust.ID
	acct.Status = "active"

	return c.JSON(fiber.Map{
		"status":    "customer_created",
		"customerId": stripeCust.ID,
		"login":     req.Login,
		"email":     req.Email,
		"plan":      req.Plan,
	})
}

// getBillingAccount returns the billing status for a GitHub account.
func getBillingAccount(c *fiber.Ctx) error {
	login := c.Params("login")
	acct, exists := accounts[login]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "account not found"})
	}
	return c.JSON(fiber.Map{
		"login":     acct.Login,
		"status":    acct.Status,
		"plan_name": acct.PlanName,
		"stripe_id": acct.StripeCustomerID,
		"github_id": acct.GitHubInstallationID,
		"created_at": acct.CreatedAt,
	})
}

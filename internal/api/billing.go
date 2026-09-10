package api

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
)

// Plan maps plan names to Stripe product IDs (set in .env).
var PlanProducts map[string]string

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	PlanProducts = map[string]string{
		"starter": os.Getenv("STRIPE_PRODUCT_STARTER"),
		"growth":  os.Getenv("STRIPE_PRODUCT_GROWTH"),
		"scale":   os.Getenv("STRIPE_PRODUCT_SCALE"),
	}
}

// RegisterBilling wires the Stripe billing endpoints.
func RegisterBilling(app *fiber.App) {
	app.Post("/api/billing/checkout", createCheckout)
	app.Post("/api/billing/customer", createCustomer)
	app.Get("/api/billing/account/:login", getBillingAccount)
	app.Get("/api/billing/plans", listPlans)
}

// listPlans returns available plans with pricing.
func listPlans(c *fiber.Ctx) error {
	type planInfo struct {
		Name      string `json:"name"`
		Price     int64  `json:"price"`
		Interval  string `json:"interval"`
		ProductID string `json:"productId"`
	}

	var plans []planInfo
	for name, productID := range PlanProducts {
		if productID == "" {
			continue
		}
		pi := planInfo{Name: name, ProductID: productID}

		if stripe.Key != "" {
			// Look up the price for this product
			params := &stripe.PriceListParams{
				Product: stripe.String(productID),
			}
			params.Filters.AddFilter("limit", "", "1")
			i := price.List(params)
			if i.Next() {
				p := i.Price()
				pi.Price = p.UnitAmount
				if p.Recurring != nil {
					pi.Interval = string(p.Recurring.Interval)
				}
			}
		}

		plans = append(plans, pi)
	}
	return c.JSON(plans)
}

// createCheckout creates a Stripe Checkout Session for a subscription.
func createCheckout(c *fiber.Ctx) error {
	if stripe.Key == "" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "billing not configured (STRIPE_SECRET_KEY missing)"})
	}

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

	productID, ok := PlanProducts[req.Plan]
	if !ok || productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("unknown plan: %s (available: starter, growth, scale)", req.Plan)})
	}

	// Look up price ID from product
	priceID, err := getPriceForProduct(productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to look up price for plan"})
	}

	acct, err := getAccount(req.Login)
	if err != nil || acct == nil {
		acct = &Account{
			Login:      req.Login,
			SourceType: "github_marketplace",
			Status:     "active",
			CreatedAt:  time.Now().UTC(),
		}
		upsertAccount(acct)
	}

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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create checkout session"})
	}

	return c.JSON(fiber.Map{
		"status":    "checkout_created",
		"sessionId": stripeSession.ID,
		"login":     req.Login,
		"plan":      req.Plan,
		"url":       stripeSession.URL,
	})
}

// createCustomer creates a Stripe customer for a marketplace account.
func createCustomer(c *fiber.Ctx) error {
	if stripe.Key == "" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "billing not configured (STRIPE_SECRET_KEY missing)"})
	}

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

	acct, err := getAccount(req.Login)
	if err != nil || acct == nil {
		acct = &Account{
			Login:      req.Login,
			SourceType: "github_marketplace",
			CreatedAt:  time.Now().UTC(),
		}
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create Stripe customer"})
	}

	acct.StripeCustomerID = stripeCust.ID
	acct.Status = "active"
	acct.Email = req.Email
	upsertAccount(acct)

	return c.JSON(fiber.Map{
		"status":     "customer_created",
		"customerId": stripeCust.ID,
		"login":      req.Login,
		"email":      req.Email,
		"plan":       req.Plan,
	})
}

// getBillingAccount returns the billing status for a GitHub account.
func getBillingAccount(c *fiber.Ctx) error {
	login := c.Params("login")
	acct, err := getAccount(login)
	if err == sql.ErrNoRows || acct == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "account not found"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"login":      acct.Login,
		"status":     acct.Status,
		"plan_name":  acct.PlanName,
		"stripe_id":  acct.StripeCustomerID,
		"github_id":  acct.GitHubInstallationID,
		"created_at": acct.CreatedAt,
	})
}

// getPriceForProduct returns the first active recurring price ID for a product.
func getPriceForProduct(productID string) (string, error) {
	params := &stripe.PriceListParams{
		Product: stripe.String(productID),
	}
	params.Filters.AddFilter("active", "", "true")
	params.Filters.AddFilter("limit", "", "10")

	i := price.List(params)
	for i.Next() {
		p := i.Price()
		if p.Recurring != nil {
			return p.ID, nil
		}
	}

	// Fallback: create a price from the product metadata
	prod, err := product.Get(productID, nil)
	if err != nil {
		return "", fmt.Errorf("product not found: %s", productID)
	}

	return "", fmt.Errorf("no recurring price found for product %s (%s)", productID, prod.Name)
}

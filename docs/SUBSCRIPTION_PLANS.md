# Subscription Plans — Aura Amber

Billing: Stripe Checkout + Stripe Customers (`/api/billing/*`).
GitHub Marketplace purchases map to the same plans via `/github/webhooks`.

| Plan | Price (suggested) | Migrations/mo | Records / migration | Connectors | Support |
|------|-------------------|---------------|---------------------|------------|---------|
| Starter | $29/mo | 3 | 10k | WP + Woo | Community |
| Growth | $99/mo | 25 | 250k | + Shopify | Email, 48h |
| Scale | $299/mo | Unlimited* | 2M | + Magento, API | Priority, 24h |

*Fair use; throttle via BullMQ concurrency.

## Stripe wiring

`.env` keys:

```
STRIPE_SECRET_KEY=
STRIPE_WEBHOOK_SECRET=
STRIPE_PRICE_ID_PLAN_BASIC=        # Starter
STRIPE_PRICE_ID_PLAN_GROWTH=       # add
STRIPE_PRICE_ID_PLAN_SCALE=        # add
```

Endpoints:

- `POST /api/billing/checkout { login, plan, email }` -> Stripe Checkout Session
- `POST /api/billing/customer { login, email, plan }` -> Stripe Customer
- `GET /api/billing/:login` -> `{ stripe_id, plan_name, status }`

## Recommendation

- Add `plan` enum (`starter|growth|scale`) validation in `billing.go`.
- Store `stripe_price_id` per account; gate `POST /api/migrations` by monthly count.
- Add `GET /api/billing/prices` returning public plan table for the landing page.

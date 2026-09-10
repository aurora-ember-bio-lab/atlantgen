# Aura Amber

Migration SaaS: WordPress / Shopify / WooCommerce / Magento → Next.js + MedusaJS + PostgreSQL.

## Status

- ✅ `ui/next-app` — migration dashboard (Next.js 16, TypeScript, Tailwind). Dashboard, Migrations, Connectors, Settings views.
- ✅ `cmd/aura-amber-api` — Go API server (Fiber framework). Endpoints: `/api/migrations`, `/api/connectors`, `/api/schema/*`, `/github/webhooks`, `/api/billing/*`.
- ✅ `worker/` — Node.js/BullMQ migration worker with 4 connectors (WordPress, Shopify, WooCommerce, Magento).
- ✅ `atlas/schema.hcl` — Database schema definition (accounts, products, orders, customers).
- ✅ `docker-compose.yml` — Full stack: Postgres 16, Redis 7, API, Worker.
- ✅ `.github/workflows/security-scan.yml` — Weekly `npm audit` on push/PR.
- ✅ Stripe billing integration with `stripe-go/v76` SDK.
- ⏳ Atlas CLI — install and apply schema to Postgres.

## Repository

- GitHub: `cargounetcom/aura-amber-saas`
- Vercel team: `aurora-ember-cyber`

## Quick Start

### Prerequisites
- Go 1.22+
- Node.js 20+
- PostgreSQL 16
- Redis 7

### Install
```bash
# Install Atlas CLI
go install ariga.io/atlas/cmd/atlas@latest

# Install UI dependencies
cd ui/next-app && npm install

# Install Worker dependencies
cd worker && npm install --legacy-peer-deps

# Install Go dependencies
cd .. && go mod tidy
```

### Development
```bash
# Start database and Redis (requires Docker)
docker compose up -d

# Apply Atlas schema to Postgres
atlas schema apply --to file://atlas/schema.hcl --url "$DATABASE_URL" --auto-approve

# Start Go API server
go run ./cmd/aura-amber-api

# Start Node.js worker
cd worker && npm run dev

# Start Next.js dashboard
cd ui/next-app && npm run dev
```

Open `http://localhost:3000` for the dashboard.

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐     ┌───────────┐
│  Next.js UI │────▶│ Go API (Fiber)│────▶│ Node Worker  │────▶│ PostgreSQL│
│  (Next 16)  │     │  (:8080)     │     │ (BullMQ/Redis)│     │  (Postgres)│
└─────────────┘     └──────┬───────┘     └──────────────┘     └───────────┘
                           │
                    ┌──────▼───────┐
                    │  Atlas CLI   │  (Schema orchestration)
                    └──────────────┘
                           │
                    ┌──────▼───────┐
                    │  GitHub App  │  (Billing/Marketplace)
                    └──────────────┘
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/migrations` | List all migration jobs |
| POST | `/api/migrations` | Create a new migration job |
| GET | `/api/connectors` | List source connectors |
| PUT | `/api/connectors/:id` | Update connector state |
| POST | `/api/schema/diff` | Diff target DB against schema |
| POST | `/api/schema/apply` | Apply schema to target DB |
| POST | `/github/webhooks` | GitHub App webhook handler |
| POST | `/api/billing/checkout` | Create Stripe Checkout Session |
| POST | `/api/billing/customer` | Create Stripe Customer |
| GET | `/api/billing/:login` | Get billing account status |

## Connectors

| Source | Status | Auth |
|--------|--------|------|
| WordPress | ✅ Connected | REST API + Application Password |
| WooCommerce | ✅ Connected | WooCommerce REST API (consumer key/secret) |
| Shopify | ⚠️ Not connected | Admin REST API (access token) |
| Magento | ⚠️ Error | REST API (Bearer token) |

## Security

```bash
cd ui/next-app
npm run scan
```

Runs `npm audit` at high severity threshold. Also runs automatically every Monday via `.github/workflows/security-scan.yml`.

## Deploy

### Vercel (UI)
```bash
vercel --prod
```

### Docker (Full Stack)
```bash
docker compose up -d
```

## Environment Variables

See `.env.example` for all required variables:
- `DATABASE_URL` — Postgres connection
- `REDIS_HOST`, `REDIS_PORT` — BullMQ queue
- `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET` — Billing
- `GITHUB_APP_ID`, `GITHUB_APP_PRIVATE_KEY`, `GITHUB_WEBHOOK_SECRET` — GitHub App
- `NEXT_PUBLIC_API_URL` — Frontend API URL
- `WORKER_URL`, `WORKER_HTTP_PORT` — Worker shim

## License

BSL 1.1 (Business Source License) — converts to MIT on 2030-09-10. See [LICENSE](./LICENSE).

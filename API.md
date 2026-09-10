# Aura Amber API Documentation

Base URL: `http://localhost:8080`

## Health

```
GET /health
```

Returns `{"status": "ok"}`.

## Migrations

### List Migrations
```
GET /api/migrations
```

Returns array of all migration jobs.

### Create Migration
```
POST /api/migrations
Content-Type: application/json

{
  "name": "my-store",
  "sourceType": "wordpress",
  "sourceConnectionInfo": {
    "baseUrl": "https://example.com",
    "username": "admin",
    "appPassword": "..."
  },
  "targetDbUrl": "postgres://aura:aura@db:5432/aura?sslmode=disable"
}
```

Returns the created job with `id`, `status: "queued"`, and `workerJobId`.

## Connectors

### List Connectors
```
GET /api/connectors
```

Returns array of connectors with `id`, `name`, `state`, `connected`, `lastSync`.

### Update Connector
```
PUT /api/connectors/:id
Content-Type: application/json

{
  "state": "connected"
}
```

Returns the updated connector.

## Schema Orchestrator

### Diff Schema
```
POST /api/schema/diff
Content-Type: application/json

{
  "targetDbUrl": "postgres://aura:aura@db:5432/aura?sslmode=disable"
}
```

Returns the SQL diff between current schema and `atlas/schema.hcl`.

### Apply Schema
```
POST /api/schema/apply
Content-Type: application/json

{
  "targetDbUrl": "postgres://aura:aura@db:5432/aura?sslmode=disable"
}
```

Applies the schema defined in `atlas/schema.hcl` using Atlas CLI.

## GitHub App Webhooks

### Webhook Endpoint
```
POST /github/webhooks
X-GitHub-Event: marketplace_purchase | installation
X-Hub-Signature-256: sha256=...
```

Handles marketplace purchase events (billing) and installation lifecycle events.

## Billing (Stripe)

### Create Checkout Session
```
POST /api/billing/checkout
Content-Type: application/json

{
  "login": "customer-name",
  "plan": "basic",
  "email": "customer@example.com"
}
```

Returns `{"status": "checkout_created", "sessionId": "..."}`.

### Create Customer
```
POST /api/billing/customer
Content-Type: application/json

{
  "login": "customer-name",
  "email": "customer@example.com",
  "plan": "basic"
}
```

Returns `{"status": "customer_created", "customerId": "..."}`.

### Get Billing Account
```
GET /api/billing/:login
```

Returns account billing status including `stripe_id`, `plan_name`, `status`.

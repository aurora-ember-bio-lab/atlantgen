# Contributing to Aura Amber

## Project Structure

```
aura-amber-saas/
├── cmd/aura-amber-api/     # Go API server entrypoint (Fiber)
├── cmd/worker/             # Go worker placeholder (worker is Node.js)
├── internal/api/           # Go API handlers (migrations, schema, billing, connectors, webhooks)
├── internal/db/            # Postgres connection
├── internal/migration/     # Atlas schema orchestration, job tracking
├── ui/next-app/            # Next.js 16 dashboard (React, Tailwind)
│   ├── app/               # App Router pages (Dashboard, Migrations, Connectors, Settings)
│   ├── components/        # Shared components (MigrationTable, Sidebar, StatCard)
│   └── lib/               # API client, type definitions
├── worker/                 # Node.js BullMQ migration worker
│   ├── src/connectors/     # Platform extractors (WordPress, Shopify, WooCommerce, Magento)
│   ├── src/loaders/        # Data normalizers and Postgres uploaders
│   └── src/server.ts       # HTTP shim for Go API to enqueue jobs
├── atlas/                  # Atlas schema definition and migration scripts
├── github-app/             # GitHub App manifest configuration
└── docker-compose.yml      # Full stack orchestration
```

## Development Workflow

1. **Fork** the repository
2. **Create a branch** from `main`
3. **Install dependencies**: `npm install` (UI + Worker), `go mod tidy` (Go)
4. **Make changes** following the code style
5. **Test**: `go test ./...`, `npx tsc --noEmit`
6. **Run linter**: `npm run lint` (UI), `go vet ./...`
7. **Commit** with a descriptive message
8. **Push** and open a Pull Request

## Adding a New Connector

1. Add connector file to `worker/src/connectors/{name}.ts`
2. Implement `extractConnectionInfo()` function
3. Add entry to `worker/src/index.ts` `extractors` map
4. Add connector to `internal/api/connectors.go` `connectors` slice
5. Add schema to `atlas/schema.hcl` if needed
6. Update `worker/src/loaders/normalize.ts` to handle the new source type
7. Update `ui/lib/connector-types.ts` with connector metadata

## Running Tests

```bash
# Go tests
go test ./...

# TypeScript linting
cd ui/next-app && npm run lint

# TypeScript type checking
npx tsc --noEmit
```

## Code Style

- **Go**: Standard Go formatting (`gofmt`)
- **TypeScript**: Strict mode enabled, no `any` types
- **SQL**: Use parameterized queries to prevent injection
- **Docker**: Multi-stage builds, alpine images

# Aura Amber — Local Test / CLI Matrix

Verified 2026-09-10 on Windows 11 + Docker Desktop 29.7.2.

## 1. Prerequisites

```powershell
go version        # 1.22+
node -v           # 20+
docker --version  # 29.x
docker compose version
```

## 2. Install (fresh clone)

```powershell
git clone https://github.com/cargounetcom/aura-amber-saas.git
cd aura-amber-saas
Copy-Item .env.example .env -Force
go mod tidy
cd worker; npm install --legacy-peer-deps; cd ..
cd ui/next-app; npm install; cd ../..
```

## 3. Start full stack (Docker — recommended)

```powershell
docker compose up -d --build
docker compose ps
# expected: aura-api (:8080), aura-db (:5432), aura-redis (:6379), aura-worker (:8090)
```

DB name is `aura_amber`. If you upgraded from an old volume with only `aura` DB:

```powershell
docker exec aura-db psql -U aura -d aura -c "CREATE DATABASE aura_amber;"
```

Tables (`accounts, customers, migrations, orders, products`) are applied via
`atlas/schema.hcl` — currently applied manually; API image now ships
`./atlas` so this also works:

```powershell
docker exec aura-api atlas schema apply --to file:///app/atlas/schema.hcl --url "postgres://aura:aura@db:5432/aura_amber?sslmode=disable" --auto-approve
```

## 4. Health checks

```powershell
curl.exe -s http://localhost:8080/health          # {"status":"ok"}
curl.exe -s http://localhost:8090/health          # {"status":"ok"}
curl.exe -s http://localhost:8080/api/migrations  # []
curl.exe -s http://localhost:8080/api/connectors  # 4 connectors
docker exec aura-db psql -U aura -d aura_amber -c "\dt"
```

## 5. E2E migration (PowerShell-safe — use Invoke-RestMethod, not curl quoting)

```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/migrations -Method Post -ContentType 'application/json' -Body '{"name":"e2e-wordpress","sourceType":"wordpress","sourceConnectionInfo":{},"targetDbUrl":"postgres://aura:aura@db:5432/aura_amber?sslmode=disable"}'
# -> status running + workerJobId

Invoke-RestMethod -Uri http://localhost:8090/enqueue -Method Post -ContentType 'application/json' -Body '{"name":"test","sourceType":"wordpress","sourceConnectionInfo":{},"targetDbUrl":"postgres://aura:aura@db:5432/aura_amber?sslmode=disable"}'
# -> {"jobId":"N"}
```

Why: `curl.exe -d "{\"a\":1}"` mangles JSON under PowerShell. `Invoke-RestMethod` is reliable.

## 6. Code checks

```powershell
go vet ./...
go test ./...            # currently: no test files — see recommendation
docker exec aura-api atlas version
cd ui/next-app; npm run lint; npm run build
cd ../../worker; npm run build
```

## 7. Local dev without Docker

```powershell
# terminal 1: postgres + redis must already run locally
go run ./cmd/aura-amber-api          # :8080, WORKER_URL=http://localhost:8090
# terminal 2:
cd worker; npm run dev              # :8090
# terminal 3:
cd ui/next-app; npm run dev         # :3000, NEXT_PUBLIC_API_URL=http://localhost:8080
```

Inside Docker the API must use `WORKER_URL=http://worker:8090`
(localhost inside a container points at itself). This is now set in
`docker-compose.yml`.

## 8. Known issues fixed 2026-09-10

- `worker: Cannot find module './queue.js'` — fixed by `"type":"module"` + `tsx` dev + `tsc build` + `node dist` start.
- `aura-worker` never listening on :8090 — fixed by `EXPOSE 8090` + `ports: 8090:8090`.
- `POST /api/migrations -> worker unreachable at http://localhost:8090` — fixed by `WORKER_URL=http://worker:8090` in compose.
- `npm run build` failed: `useState` without `"use client"` — fixed in `app/migrations/page.tsx`.
- `npm run build` fetched live API at build time — fixed with `force-dynamic` + try/catch fallback in `app/page.tsx`, `app/connectors/page.tsx`.
- `POSTGRES_DB=aura` vs `DATABASE_URL=.../aura_amber` mismatch — fixed to `aura_amber` everywhere.
- API image missing `./atlas` — fixed by `COPY --from=build /app/atlas ./atlas`.
- Port :8080 double-bind (`aura-amber-api.exe` + Docker) — kill local exe before `docker compose up`.

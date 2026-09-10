# Release Process — aura-amber-saas

Repo: https://github.com/cargounetcom/aura-amber-saas
Default branch: `main`

## Versioning

SemVer `vMAJOR.MINOR.PATCH` (e.g. `v0.2.0`).

## Pre-release checklist

1. `go vet ./...` clean
2. `cd ui/next-app; npm run lint; npm run build` green
3. `cd worker; npm run build` green
4. `docker compose up -d --build` + health checks in `docs/LOCAL_TEST.md`
5. E2E `POST /api/migrations` returns `running` + `workerJobId`
6. Update `CHANGELOG.md` (create if missing)

## Cut a release

```powershell
git checkout main
git pull origin main
git tag -a v0.2.0 -m "v0.2.0: worker ESM fix, compose networking, UI build fixes"
git push origin v0.2.0
gh release create v0.2.0 --title "v0.2.0" --notes-file CHANGELOG.md
```

Suggested `v0.2.0` notes:

- Fix worker boot (`type:module`, tsx build/start, EXPOSE 8090)
- Fix API->worker networking (`WORKER_URL=http://worker:8090`)
- Fix UI build (`use client`, force-dynamic fallbacks)
- Unify DB name to `aura_amber`, ship `./atlas` in API image
- Docs: LOCAL_TEST, VERCEL_DEPLOY, SUBSCRIPTION_PLANS, MARKETING

## Post-release

- Vercel auto-deploys `ui/next-app` on tag if connected (see `docs/VERCEL_DEPLOY.md`)
- Verify production `/api/migrations`, `/api/connectors` via `NEXT_PUBLIC_API_URL`
- Open GitHub Release, attach Docker images if needed

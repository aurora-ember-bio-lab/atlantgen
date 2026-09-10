# Vercel Deploy — UI front (`ui/next-app`)

Team: `aurora-ember-cyber` (per README).

## Fresh deploy (new project)

```powershell
cd ui/next-app
npm i -g vercel
vercel link          # team aurora-ember-cyber, project aura-amber
vercel --prod
```

Vercel settings:

- Root Directory: `ui/next-app`
- Framework: Next.js
- Build: `npm run build`
- Env: `NEXT_PUBLIC_API_URL=https://<your-api>/` (e.g. Render/Fly for Go API)

`ui/next-app/vercel.json` pins the framework and denies `/api/*` rewrites
(the UI has no Route Handlers; all data comes from the Go API).

## Landing + legal routes (in this deploy)

- `/` dashboard (requires API, falls back to empty state at build)
- `/migrations`, `/connectors`, `/settings`
- `/privacy`, `/cookies`, `/disclaimer`, `/ads-disclosure`

Link these in footers/ads for compliance.

## Preview flow

Every PR gets a Preview URL. Test:

```
<preview>/migrations
<preview>/privacy
```

Promote to Production after `npm run build` passes.

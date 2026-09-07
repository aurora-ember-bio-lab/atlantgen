# Aura Amber

Migration SaaS: WordPress / Shopify / WooCommerce / Magento → Next.js + MedusaJS + PostgreSQL.

## Status

- ✅ `ui/next-app` — migration dashboard (Next.js 16, TypeScript, Tailwind). Dashboard, Migrations, Connectors, Settings views.
- ⏳ Backend (Go API, migration worker, Atlas schema orchestrator, GitHub App billing) — not started yet.

## Repository

- GitHub: `cargounetcom/aura-amber-saas`
- Vercel team: `aurora-ember-cyber`

## Development

```bash
cd ui/next-app
npm install
npm run dev
```

Open `http://localhost:3000`.

## Security

```bash
cd ui/next-app
npm run scan
```

Runs `npm audit` at high severity threshold. Also runs automatically every Monday via `.github/workflows/security-scan.yml`.

## Deploy

```bash
vercel --prod
```

## License

MIT — see [LICENSE](./LICENSE).

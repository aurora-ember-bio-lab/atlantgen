# Marketing — Aura Amber

Tagline: **From legacy store to modern stack in hours, not weeks.**
Sources: WordPress / WooCommerce / Shopify / Magento.
Target: Next.js + MedusaJS + PostgreSQL.

## Positioning

- Pain: replatforming is manual, risky, SEO-breaking.
- Promise: connectors + normalizers + uploaders + Atlas schema + dashboard.
- Proof: live dashboard, BullMQ progress, idempotent upserts.

## Channels

1. GitHub repo (primary): polished README, topics
   `migration, wordpress, shopify, woocommerce, magento, nextjs, medusajs, postgres`.
2. Vercel deployment URL on every release + landing page.
3. Short demos: WP -> Postgres products/orders/customers in <5 min.
4. SEO: `/privacy`, `/cookies`, `/disclaimer`, `/ads-disclosure` pages (trust + ad compliance).

## Copy blocks

- Hero: "Migrate your store to Next.js + MedusaJS without losing orders, customers, or SEO."
- CTA: "Start a migration" -> `/migrations`.
- Trust: "BSL 1.1 licensed (converts to MIT 2030). Your data stays in your Postgres."

## GitHub hygiene

```powershell
gh repo edit cargounetcom/aura-amber-saas --description "WordPress/Shopify/Woo/Magento to Next.js+MedusaJS+Postgres migration SaaS" --add-topic migration --add-topic nextjs --add-topic medusajs --add-topic postgres
```

Pin `docs/LOCAL_TEST.md` and the Vercel link in README.

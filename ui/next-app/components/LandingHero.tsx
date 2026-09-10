"use client";

import { useState } from "react";
import GlassTabs from "./GlassTabs";

const connectors = [
  { id: "wordpress", name: "WordPress", icon: "W", color: "bg-blue-500/10 text-blue-400 border-blue-500/30" },
  { id: "woocommerce", name: "WooCommerce", icon: "Wc", color: "bg-purple-500/10 text-purple-400 border-purple-500/30" },
  { id: "shopify", name: "Shopify", icon: "S", color: "bg-green-500/10 text-green-400 border-green-500/30" },
  { id: "magento", name: "Magento", icon: "M", color: "bg-orange-500/10 text-orange-400 border-orange-500/30" },
];

const steps = [
  {
    num: "01",
    title: "Connect your store",
    desc: "Point Aura Amber at your WordPress, Shopify, WooCommerce, or Magento instance. No plugins. No downtime.",
    code: `POST /api/connectors/wordpress
{
  "baseUrl": "https://my-store.com",
  "username": "admin",
  "appPassword": "xxxx xxxx xxxx xxxx"
}`,
    lang: "bash",
  },
  {
    num: "02",
    title: "Start migration",
    desc: "One command kicks off extraction, normalization, and upload. BullMQ handles retries and progress.",
    code: `POST /api/migrations
{
  "name": "my-store",
  "sourceType": "wordpress",
  "targetDbUrl": "postgres://you:pass@host/your_db"
}`,
    lang: "json",
  },
  {
    num: "03",
    title: "Go live on Next.js + MedusaJS",
    desc: "Products, orders, customers land in PostgreSQL with Atlas-managed schema. Deploy to Vercel in seconds.",
    code: `vercel --prod
# https://your-store.vercel.app ✓`,
    lang: "bash",
  },
];

const features = [
  {
    title: "Declarative schema",
    desc: "Define your target schema in HCL. Atlas computes the safe migration plan automatically.",
    before: `-- Old way: hand-write every ALTER TABLE
ALTER TABLE products ADD COLUMN price_cents INT;
ALTER TABLE products ADD COLUMN source_id TEXT;
CREATE UNIQUE INDEX products_source_idx ON products(source_type, source_id);`,
    after: `-- Atlas way: define desired state, it computes the diff
table "products" {
  column "id"          { type = uuid }
  column "title"       { type = text }
  column "price_cents" { type = int }
  column "source_id"   { type = text }
  index "products_source_idx" {
    columns = [column.source_type, column.source_id]
    unique  = true
  }
}`,
    langBefore: "sql",
    langAfter: "hcl",
  },
  {
    title: "Real-time progress",
    desc: "BullMQ job queue with Redis. Watch extraction, normalization, and upload live from the dashboard.",
    before: `# Old way: scripts that silently fail
python migrate.py --source=wordpress 2>&1 | tee log.txt
# Did it work? Who knows.`,
    after: `# Aura Amber: structured progress
GET /jobs/42
{
  "id": "42",
  "state": "active",
  "progress": {
    "extract": "100%",
    "normalize": "75%",
    "upload": "0%"
  }
}`,
    langBefore: "bash",
    langAfter: "json",
  },
  {
    title: "Idempotent upserts",
    desc: "Run migrations as many times as you want. Products, orders, and customers are upserted, never duplicated.",
    before: `-- Old way: INSERT or crash on duplicate
INSERT INTO products (title, price) VALUES ('Widget', 2990);
-- ERROR: duplicate key value violates unique constraint`,
    after: `-- Aura Amber: safe upsert
INSERT INTO products (id, title, source_type, source_id, price_cents)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (source_type, source_id) DO UPDATE
SET title = EXCLUDED.title, price_cents = EXCLUDED.price_cents;`,
    langBefore: "sql",
    langAfter: "sql",
  },
];

const stats = [
  { value: "4", label: "Source platforms" },
  { value: "<5min", label: "First migration" },
  { value: "100%", label: "Idempotent" },
  { value: "0", label: "Downtime" },
];

export default function LandingHero() {
  const [activeStep, setActiveStep] = useState(0);
  const [activeFeature, setActiveFeature] = useState(0);
  const [showBefore, setShowBefore] = useState(false);

  return (
    <div className="min-h-screen">
      {/* ── Nav ─────────────────────────────────────────────────── */}
      <nav className="sticky top-0 z-50 border-b border-neutral-800/60 bg-neutral-950/80 backdrop-blur-xl">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-2.5">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-amber-500 text-sm font-bold text-neutral-950">
              A
            </div>
            <span className="text-lg font-semibold tracking-tight">Aura Amber</span>
          </div>
          <div className="hidden items-center gap-8 md:flex">
            <a href="#features" className="text-sm text-neutral-400 transition-colors hover:text-amber-400">Features</a>
            <a href="#how-it-works" className="text-sm text-neutral-400 transition-colors hover:text-amber-400">How it works</a>
            <a href="#schema" className="text-sm text-neutral-400 transition-colors hover:text-amber-400">Schema</a>
            <a href="/docs" className="text-sm text-neutral-400 transition-colors hover:text-amber-400">Docs</a>
          </div>
          <div className="flex items-center gap-3">
            <a href="/dashboard" className="rounded-lg border border-neutral-700 px-4 py-2 text-sm font-medium text-neutral-300 transition-colors hover:border-amber-500 hover:text-amber-400">
              Dashboard
            </a>
            <a href="#get-started" className="rounded-lg bg-amber-500 px-4 py-2 text-sm font-semibold text-neutral-950 transition-colors hover:bg-amber-400">
              Get Started
            </a>
          </div>
        </div>
      </nav>

      {/* ── Hero ────────────────────────────────────────────────── */}
      <section className="relative overflow-hidden pb-20 pt-24">
        {/* gradient glow */}
        <div className="pointer-events-none absolute left-1/2 top-0 h-[600px] w-[800px] -translate-x-1/2 rounded-full bg-amber-500/8 blur-[120px]" />

        <div className="relative mx-auto max-w-4xl px-6 text-center">
          <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-amber-500/30 bg-amber-500/10 px-4 py-1.5 text-xs font-medium text-amber-400">
            <span className="h-1.5 w-1.5 rounded-full bg-amber-400 animate-pulse" />
            Now with Atlas schema orchestration
          </div>

          <h1 className="text-5xl font-bold leading-tight tracking-tight md:text-7xl">
            Migrate your store
            <br />
            <span className="bg-gradient-to-r from-amber-400 via-amber-300 to-orange-400 bg-clip-text text-transparent">
              without losing a beat
            </span>
          </h1>

          <p className="mx-auto mt-6 max-w-2xl text-lg text-neutral-400">
            WordPress, Shopify, WooCommerce, Magento → Next.js + MedusaJS + PostgreSQL.
            One command. Zero downtime. Products, orders, customers — all there.
          </p>

          <div className="mt-10 flex flex-wrap items-center justify-center gap-4">
            <a href="#get-started" className="rounded-xl bg-amber-500 px-8 py-3.5 text-base font-semibold text-neutral-950 shadow-lg shadow-amber-500/20 transition-all hover:bg-amber-400 hover:shadow-amber-500/30">
              Start migrating →
            </a>
            <a href="#how-it-works" className="rounded-xl border border-neutral-700 px-8 py-3.5 text-base font-medium text-neutral-300 transition-colors hover:border-amber-500 hover:text-amber-400">
              See how it works
            </a>
          </div>

          {/* connector badges */}
          <div className="mt-12 flex flex-wrap items-center justify-center gap-3">
            {connectors.map((c) => (
              <div key={c.id} className={`inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium ${c.color}`}>
                <span className="flex h-6 w-6 items-center justify-center rounded-md bg-neutral-800 text-xs font-bold">{c.icon}</span>
                {c.name}
              </div>
            ))}
            <span className="text-neutral-600">→</span>
            <div className="inline-flex items-center gap-2 rounded-full border border-neutral-700 bg-neutral-900/60 px-4 py-2 text-sm font-medium text-neutral-300">
              <span className="flex h-6 w-6 items-center justify-center rounded-md bg-neutral-800 text-xs font-bold text-amber-400">N</span>
              Next.js + MedusaJS
            </div>
          </div>
        </div>
      </section>

      {/* ── Stats bar ───────────────────────────────────────────── */}
      <section className="border-y border-neutral-800/60 bg-neutral-900/30">
        <div className="mx-auto grid max-w-5xl grid-cols-2 gap-8 px-6 py-10 md:grid-cols-4">
          {stats.map((s) => (
            <div key={s.label} className="text-center">
              <div className="text-3xl font-bold text-amber-400">{s.value}</div>
              <div className="mt-1 text-sm text-neutral-500">{s.label}</div>
            </div>
          ))}
        </div>
      </section>

      {/* ── How it works (step tabs) ────────────────────────────── */}
      <section id="how-it-works" className="py-24">
        <div className="mx-auto max-w-6xl px-6">
          <div className="mb-12 text-center">
            <h2 className="text-3xl font-bold md:text-4xl">Three steps. That&apos;s it.</h2>
            <p className="mt-3 text-neutral-400">Connect, migrate, deploy. No scripts to maintain.</p>
          </div>

          <div className="grid gap-8 lg:grid-cols-2">
            {/* step selector */}
            <div className="space-y-3">
              {steps.map((s, i) => (
                <button
                  key={s.num}
                  onClick={() => setActiveStep(i)}
                  className={`w-full rounded-xl border p-5 text-left transition-all ${
                    activeStep === i
                      ? "border-amber-500/50 bg-amber-500/5 shadow-lg shadow-amber-500/5"
                      : "border-neutral-800 bg-neutral-900/30 hover:border-neutral-700"
                  }`}
                >
                  <div className="flex items-start gap-4">
                    <span className={`text-2xl font-bold ${activeStep === i ? "text-amber-400" : "text-neutral-600"}`}>
                      {s.num}
                    </span>
                    <div>
                      <h3 className="text-base font-semibold text-neutral-100">{s.title}</h3>
                      <p className="mt-1 text-sm text-neutral-500">{s.desc}</p>
                    </div>
                  </div>
                </button>
              ))}
            </div>

            {/* code block */}
            <div className="rounded-xl border border-neutral-800 bg-neutral-950 overflow-hidden">
              <div className="flex items-center gap-2 border-b border-neutral-800 px-4 py-3">
                <div className="h-3 w-3 rounded-full bg-red-500/60" />
                <div className="h-3 w-3 rounded-full bg-yellow-500/60" />
                <div className="h-3 w-3 rounded-full bg-green-500/60" />
                <span className="ml-2 text-xs text-neutral-600">{steps[activeStep].lang}</span>
              </div>
              <pre className="p-5 text-sm leading-relaxed text-neutral-300 overflow-x-auto">
                <code>{steps[activeStep].code}</code>
              </pre>
            </div>
          </div>
        </div>
      </section>

      {/* ── Features (before/after tabs) ────────────────────────── */}
      <section id="features" className="border-t border-neutral-800/60 bg-neutral-900/20 py-24">
        <div className="mx-auto max-w-6xl px-6">
          <div className="mb-12 text-center">
            <h2 className="text-3xl font-bold md:text-4xl">
              Goodbye manual migrations.
              <br />
              <span className="text-amber-400">Hello Aura Amber.</span>
            </h2>
            <p className="mt-3 text-neutral-400">Declarative, idempotent, real-time. The modern way to replatform.</p>
          </div>

          {/* feature tabs */}
          <div className="mb-8 flex flex-wrap justify-center gap-2">
            {features.map((f, i) => (
              <button
                key={f.title}
                onClick={() => { setActiveFeature(i); setShowBefore(false); }}
                className={`rounded-lg px-5 py-2.5 text-sm font-medium transition-all ${
                  activeFeature === i
                    ? "bg-amber-500/10 text-amber-400 border border-amber-500/30"
                    : "border border-neutral-800 text-neutral-500 hover:text-neutral-300"
                }`}
              >
                {f.title}
              </button>
            ))}
          </div>

          {/* active feature */}
          <div className="rounded-2xl border border-neutral-800 bg-neutral-950/80 p-8">
            <div className="mb-6 flex items-center justify-between">
              <div>
                <h3 className="text-xl font-semibold">{features[activeFeature].title}</h3>
                <p className="mt-1 text-sm text-neutral-400">{features[activeFeature].desc}</p>
              </div>
              <div className="flex rounded-lg border border-neutral-800 bg-neutral-900 p-1">
                <button
                  onClick={() => setShowBefore(false)}
                  className={`rounded-md px-3 py-1.5 text-xs font-medium transition-all ${!showBefore ? "bg-amber-500/10 text-amber-400" : "text-neutral-500"}`}
                >
                  Aura Amber
                </button>
                <button
                  onClick={() => setShowBefore(true)}
                  className={`rounded-md px-3 py-1.5 text-xs font-medium transition-all ${showBefore ? "bg-red-500/10 text-red-400" : "text-neutral-500"}`}
                >
                  Old way
                </button>
              </div>
            </div>

            <div className="rounded-xl border border-neutral-800 bg-neutral-950 overflow-hidden">
              <div className="flex items-center gap-2 border-b border-neutral-800 px-4 py-3">
                <div className="h-3 w-3 rounded-full bg-red-500/60" />
                <div className="h-3 w-3 rounded-full bg-yellow-500/60" />
                <div className="h-3 w-3 rounded-full bg-green-500/60" />
                <span className="ml-2 text-xs text-neutral-600">
                  {showBefore ? features[activeFeature].langBefore : features[activeFeature].langAfter}
                </span>
              </div>
              <pre className="p-5 text-sm leading-relaxed text-neutral-300 overflow-x-auto">
                <code>{showBefore ? features[activeFeature].before : features[activeFeature].after}</code>
              </pre>
            </div>
          </div>
        </div>
      </section>

      {/* ── Schema as Code ──────────────────────────────────────── */}
      <section id="schema" className="py-24">
        <div className="mx-auto max-w-6xl px-6">
          <div className="grid gap-12 lg:grid-cols-2">
            <div>
              <h2 className="text-3xl font-bold md:text-4xl">
                Schema as Code.
                <br />
                <span className="text-amber-400">Like Terraform, for your store.</span>
              </h2>
              <p className="mt-4 text-neutral-400">
                Define your PostgreSQL schema in HCL. Atlas plans, lints, and deploys safe migrations.
                Products, orders, customers — all version-controlled.
              </p>
              <div className="mt-8 space-y-4">
                {[
                  "Declarative or versioned — your choice",
                  "Destructive change detection before deploy",
                  "Automatic diff between desired and live state",
                  "CI/CD ready: GitHub Actions, GitLab CI, Terraform",
                ].map((item) => (
                  <div key={item} className="flex items-start gap-3">
                    <div className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-amber-500/10">
                      <svg className="h-3 w-3 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                      </svg>
                    </div>
                    <span className="text-sm text-neutral-300">{item}</span>
                  </div>
                ))}
              </div>
            </div>
            <div className="rounded-xl border border-neutral-800 bg-neutral-950 overflow-hidden">
              <div className="flex items-center gap-2 border-b border-neutral-800 px-4 py-3">
                <div className="h-3 w-3 rounded-full bg-red-500/60" />
                <div className="h-3 w-3 rounded-full bg-yellow-500/60" />
                <div className="h-3 w-3 rounded-full bg-green-500/60" />
                <span className="ml-2 text-xs text-neutral-600">atlas/schema.hcl</span>
              </div>
              <pre className="p-5 text-sm leading-relaxed text-neutral-300 overflow-x-auto">
                <code>{`schema "public" {
  charset = "utf8"
}

table "products" {
  column "id" { type = uuid }
  column "title" { type = text }
  column "source_type" { type = text }
  column "source_id" { type = text }
  column "price_cents" { type = int }

  primary_key {
    columns = [column.id]
  }

  index "products_source_idx" {
    columns = [column.source_type, column.source_id]
    unique  = true
  }
}

table "orders" {
  column "id" { type = uuid }
  column "total_cents" { type = int }
  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }
}`}</code>
              </pre>
            </div>
          </div>
        </div>
      </section>

      {/* ── Architecture ────────────────────────────────────────── */}
      <section className="border-t border-neutral-800/60 bg-neutral-900/20 py-24">
        <div className="mx-auto max-w-4xl px-6 text-center">
          <h2 className="text-3xl font-bold md:text-4xl">Built on tools you already trust</h2>
          <p className="mt-3 text-neutral-400">Go, Fiber, BullMQ, PostgreSQL, Redis, Atlas, Next.js, Vercel.</p>

          <div className="mt-12 grid grid-cols-2 gap-4 md:grid-cols-4">
            {[
              { name: "Go + Fiber", desc: "API server" },
              { name: "BullMQ + Redis", desc: "Job queue" },
              { name: "PostgreSQL", desc: "Your data" },
              { name: "Atlas CLI", desc: "Schema mgmt" },
              { name: "Next.js 16", desc: "Dashboard" },
              { name: "Stripe", desc: "Billing" },
              { name: "GitHub App", desc: "Marketplace" },
              { name: "Vercel", desc: "Frontend deploy" },
            ].map((t) => (
              <div key={t.name} className="rounded-xl border border-neutral-800 bg-neutral-950/60 p-4">
                <div className="text-sm font-semibold text-neutral-100">{t.name}</div>
                <div className="mt-1 text-xs text-neutral-500">{t.desc}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Testimonials ────────────────────────────────────────── */}
      <section className="py-24">
        <div className="mx-auto max-w-5xl px-6">
          <h2 className="mb-12 text-center text-3xl font-bold md:text-4xl">What developers say</h2>
          <div className="grid gap-6 md:grid-cols-3">
            {[
              {
                quote: "We migrated 12k products and 80k orders from WooCommerce in under 4 minutes. Zero data loss.",
                name: "Migration lead",
                role: "E-commerce team",
              },
              {
                quote: "The Atlas schema integration means we never have to think about migrations again. It just works.",
                name: "Backend engineer",
                role: "SaaS platform",
              },
              {
                quote: "Finally a replatforming tool that doesn't require a month of downtime. This is the future.",
                name: "CTO",
                role: "DTC brand",
              },
            ].map((t) => (
              <div key={t.name} className="rounded-xl border border-neutral-800 bg-neutral-900/40 p-6">
                <div className="mb-4 text-amber-400">★★★★★</div>
                <p className="text-sm text-neutral-300 leading-relaxed">&ldquo;{t.quote}&rdquo;</p>
                <div className="mt-4 border-t border-neutral-800 pt-4">
                  <div className="text-sm font-medium text-neutral-100">{t.name}</div>
                  <div className="text-xs text-neutral-500">{t.role}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Glass Tabs (Plans, Legal, License) ─────────────────── */}
      <GlassTabs />

      {/* ── CTA ─────────────────────────────────────────────────── */}
      <section id="get-started" className="border-t border-neutral-800/60 py-24">
        <div className="mx-auto max-w-3xl px-6 text-center">
          <h2 className="text-4xl font-bold md:text-5xl">
            Ready to migrate?
          </h2>
          <p className="mt-4 text-lg text-neutral-400">
            Clone the repo, run <code className="rounded bg-neutral-800 px-2 py-0.5 text-amber-400">docker compose up -d</code>, and start your first migration in under 5 minutes.
          </p>

          <div className="mt-8 rounded-xl border border-neutral-800 bg-neutral-950 p-6 text-left">
            <div className="flex items-center gap-2 mb-3">
              <div className="h-3 w-3 rounded-full bg-red-500/60" />
              <div className="h-3 w-3 rounded-full bg-yellow-500/60" />
              <div className="h-3 w-3 rounded-full bg-green-500/60" />
              <span className="ml-2 text-xs text-neutral-600">terminal</span>
            </div>
            <pre className="text-sm leading-relaxed text-neutral-300"><code>{`git clone https://github.com/cargounetcom/aura-amber-saas.git
cd aura-amber-saas
docker compose up -d --build
# → API: http://localhost:8080
# → Dashboard: http://localhost:3000
# → Worker: http://localhost:8090`}</code></pre>
          </div>

          <div className="mt-8 flex flex-wrap items-center justify-center gap-4">
            <a href="https://github.com/cargounetcom/aura-amber-saas" className="rounded-xl bg-amber-500 px-8 py-3.5 text-base font-semibold text-neutral-950 shadow-lg shadow-amber-500/20 transition-all hover:bg-amber-400">
              View on GitHub →
            </a>
            <a href="/dashboard" className="rounded-xl border border-neutral-700 px-8 py-3.5 text-base font-medium text-neutral-300 transition-colors hover:border-amber-500 hover:text-amber-400">
              Open Dashboard
            </a>
          </div>
        </div>
      </section>

      {/* ── Footer ──────────────────────────────────────────────── */}
      <footer className="border-t border-neutral-800/60 bg-neutral-950">
        <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-6 px-6 py-10 md:flex-row">
          <div className="flex items-center gap-2">
            <div className="flex h-7 w-7 items-center justify-center rounded-md bg-amber-500 text-xs font-bold text-neutral-950">A</div>
            <span className="text-sm font-semibold">Aura Amber</span>
          </div>
          <div className="flex flex-wrap gap-6 text-sm text-neutral-500">
            <a href="/privacy" className="hover:text-amber-400">Privacy</a>
            <a href="/cookies" className="hover:text-amber-400">Cookies</a>
            <a href="/disclaimer" className="hover:text-amber-400">Disclaimer</a>
            <a href="/ads-disclosure" className="hover:text-amber-400">Ads</a>
            <a href="https://github.com/cargounetcom/aura-amber-saas" className="hover:text-amber-400">GitHub</a>
          </div>
          <div className="text-xs text-neutral-600">MIT License © 2026 Aurora Ember Bio Lab Ltd.</div>
        </div>
      </footer>
    </div>
  );
}

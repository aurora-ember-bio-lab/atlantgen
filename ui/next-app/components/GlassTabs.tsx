"use client";

import { useState } from "react";

const tabs = [
  { id: "plans", label: "Plans" },
  { id: "privacy", label: "Privacy" },
  { id: "license", label: "License" },
  { id: "cookies", label: "Cookies" },
  { id: "disclaimer", label: "Disclaimer" },
  { id: "ads", label: "Ads" },
];

const plans = [
  {
    name: "Starter",
    price: "$29",
    period: "/mo",
    features: ["3 migrations/mo", "10k records/migration", "WordPress + WooCommerce", "Community support"],
    cta: "Start free trial",
    highlighted: false,
  },
  {
    name: "Growth",
    price: "$99",
    period: "/mo",
    features: ["25 migrations/mo", "250k records/migration", "+ Shopify", "Email support (48h)", "Priority queue"],
    cta: "Start free trial",
    highlighted: true,
  },
  {
    name: "Scale",
    price: "$299",
    period: "/mo",
    features: ["Unlimited migrations", "2M records/migration", "+ Magento + API", "Priority support (24h)", "Custom connectors", "Dedicated instance"],
    cta: "Contact sales",
    highlighted: false,
  },
];

export default function GlassTabs() {
  const [activeTab, setActiveTab] = useState("plans");

  return (
    <section className="py-24">
      <div className="mx-auto max-w-5xl px-6">
        {/* Tab bar */}
        <div className="mb-8 flex flex-wrap justify-center gap-2">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`relative rounded-xl px-5 py-2.5 text-sm font-medium transition-all duration-300 ${
                activeTab === tab.id
                  ? "bg-amber-500/10 text-amber-400 shadow-lg shadow-amber-500/10 border border-amber-500/30"
                  : "border border-transparent text-neutral-500 hover:text-neutral-300 hover:bg-neutral-800/50"
              }`}
            >
              {tab.label}
              {activeTab === tab.id && (
                <span className="absolute inset-0 rounded-xl border border-amber-500/20 bg-amber-500/5 backdrop-blur-sm" />
              )}
            </button>
          ))}
        </div>

        {/* Tab content — glass panel */}
        <div className="relative overflow-hidden rounded-2xl border border-neutral-800/60 bg-neutral-900/40 backdrop-blur-xl shadow-2xl shadow-black/20">
          {/* subtle gradient overlay */}
          <div className="pointer-events-none absolute inset-0 bg-gradient-to-br from-amber-500/5 via-transparent to-orange-500/5" />

          <div className="relative p-8 md:p-10">
            {activeTab === "plans" && (
              <div>
                <h3 className="mb-2 text-2xl font-bold text-center">Choose your plan</h3>
                <p className="mb-10 text-center text-neutral-400">Scale as you grow. Cancel anytime.</p>
                <div className="grid gap-6 md:grid-cols-3">
                  {plans.map((plan) => (
                    <div
                      key={plan.name}
                      className={`relative rounded-xl border p-6 transition-all duration-300 ${
                        plan.highlighted
                          ? "border-amber-500/50 bg-amber-500/5 shadow-lg shadow-amber-500/10 scale-[1.02]"
                          : "border-neutral-800 bg-neutral-950/50 hover:border-neutral-700"
                      }`}
                    >
                      {plan.highlighted && (
                        <div className="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-amber-500 px-3 py-0.5 text-xs font-bold text-neutral-950">
                          MOST POPULAR
                        </div>
                      )}
                      <h4 className="text-lg font-semibold text-neutral-100">{plan.name}</h4>
                      <div className="mt-3 flex items-baseline gap-1">
                        <span className="text-4xl font-bold text-amber-400">{plan.price}</span>
                        <span className="text-sm text-neutral-500">{plan.period}</span>
                      </div>
                      <ul className="mt-6 space-y-3">
                        {plan.features.map((f) => (
                          <li key={f} className="flex items-start gap-2 text-sm text-neutral-300">
                            <svg className="mt-0.5 h-4 w-4 shrink-0 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                            </svg>
                            {f}
                          </li>
                        ))}
                      </ul>
                      <button
                        className={`mt-8 w-full rounded-lg py-2.5 text-sm font-semibold transition-all ${
                          plan.highlighted
                            ? "bg-amber-500 text-neutral-950 hover:bg-amber-400"
                            : "border border-neutral-700 text-neutral-300 hover:border-amber-500 hover:text-amber-400"
                        }`}
                      >
                        {plan.cta}
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {activeTab === "privacy" && (
              <div className="max-w-3xl mx-auto">
                <h3 className="mb-4 text-2xl font-bold">Privacy Policy</h3>
                <p className="mb-4 text-xs text-neutral-500">Last updated: 10 Sep 2026</p>
                <div className="space-y-4 text-sm text-neutral-300 leading-relaxed">
                  <p>Aura Amber migrates store data (products, orders, customers) into <strong className="text-neutral-100">your own PostgreSQL</strong>. We do not sell personal data.</p>
                  <h4 className="font-semibold text-neutral-100">Data we process</h4>
                  <ul className="list-disc pl-5 space-y-1 text-neutral-400">
                    <li>Source store records you explicitly connect</li>
                    <li>Billing contact (email, login) via Stripe</li>
                    <li>Operational logs (job status, errors)</li>
                  </ul>
                  <h4 className="font-semibold text-neutral-100">Processors</h4>
                  <p className="text-neutral-400">PostgreSQL, Redis/BullMQ, Stripe, GitHub, Vercel.</p>
                  <h4 className="font-semibold text-neutral-100">Rights</h4>
                  <p className="text-neutral-400">Request export/deletion of billing data anytime.</p>
                </div>
              </div>
            )}

            {activeTab === "license" && (
              <div className="max-w-3xl mx-auto">
                <h3 className="mb-4 text-2xl font-bold">MIT License</h3>
                <p className="mb-4 text-xs text-neutral-500">Copyright © 2026 Aurora Ember Bio Lab Ltd.</p>
                <div className="rounded-xl border border-neutral-800 bg-neutral-950/60 p-6 text-sm text-neutral-400 leading-relaxed font-mono">
                  <p>Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the &quot;Software&quot;), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software.</p>
                  <p className="mt-4">THE SOFTWARE IS PROVIDED &quot;AS IS&quot;, WITHOUT WARRANTY OF ANY KIND.</p>
                </div>
                <p className="mt-4 text-xs text-neutral-500">Full text: <a href="/license" className="text-amber-400 hover:underline">LICENSE</a></p>
              </div>
            )}

            {activeTab === "cookies" && (
              <div className="max-w-3xl mx-auto">
                <h3 className="mb-4 text-2xl font-bold">Cookie Policy</h3>
                <p className="mb-4 text-xs text-neutral-500">Last updated: 10 Sep 2026</p>
                <div className="space-y-4 text-sm text-neutral-300 leading-relaxed">
                  <p>We use only strictly-necessary cookies (session, CSRF, consent flag). No advertising trackers.</p>
                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="rounded-xl border border-neutral-800 bg-neutral-950/50 p-4">
                      <div className="text-sm font-semibold text-amber-400">Necessary</div>
                      <p className="mt-1 text-xs text-neutral-500">Auth session, cookie-consent (12 months)</p>
                    </div>
                    <div className="rounded-xl border border-neutral-800 bg-neutral-950/50 p-4">
                      <div className="text-sm font-semibold text-neutral-400">Optional</div>
                      <p className="mt-1 text-xs text-neutral-500">Analytics, ad measurement (only if enabled)</p>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {activeTab === "disclaimer" && (
              <div className="max-w-3xl mx-auto">
                <h3 className="mb-4 text-2xl font-bold">Disclaimer</h3>
                <p className="mb-4 text-xs text-neutral-500">Last updated: 10 Sep 2026</p>
                <div className="space-y-4 text-sm text-neutral-300 leading-relaxed">
                  <p>Aura Amber is provided under the MIT License <strong className="text-neutral-100">&quot;as is&quot;</strong>, without warranty.</p>
                  <ul className="list-disc pl-5 space-y-2 text-neutral-400">
                    <li>Always run a trial migration first</li>
                    <li>Back up your source store and target database</li>
                    <li>SEO, theme fidelity, and plugin parity vary by platform</li>
                    <li>Pricing is informational; Stripe Checkout price governs</li>
                  </ul>
                </div>
              </div>
            )}

            {activeTab === "ads" && (
              <div className="max-w-3xl mx-auto">
                <h3 className="mb-4 text-2xl font-bold">Ads & Affiliate Disclosure</h3>
                <p className="mb-4 text-xs text-neutral-500">Last updated: 10 Sep 2026</p>
                <div className="space-y-4 text-sm text-neutral-300 leading-relaxed">
                  <p>Aura Amber itself serves <strong className="text-neutral-100">no third-party ads</strong>.</p>
                  <p className="text-neutral-400">If marketing pages link to paid services (hosting, Stripe, Vercel), those may be affiliate or referral links. Sponsored content will be labeled &quot;Sponsored&quot;.</p>
                  <p className="text-neutral-400">Ad partners on your own storefront must be disclosed in your store&apos;s privacy/cookie pages.</p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}

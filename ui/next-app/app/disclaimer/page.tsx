export const metadata = { title: "Disclaimer — Aura Amber" };

export default function DisclaimerPage() {
  return (
    <main className="mx-auto max-w-3xl p-6 text-sm leading-6 text-neutral-300">
      <h1 className="mb-4 text-2xl font-semibold text-neutral-100">Disclaimer</h1>
      <p className="text-neutral-500">Last updated: 10 Sep 2026</p>
      <div className="mt-6 space-y-3 text-neutral-400">
        <p>Aura Amber is provided under the Business Source License 1.1 (BSL 1.1) &quot;as is&quot;, without warranty. Converts to MIT on 2030-09-10. Always run a trial migration and back up your source store and target database first.</p>
        <p>SEO, theme fidelity, and plugin parity vary by source platform and require manual review after migration.</p>
        <p>Pricing on the site is informational; the Stripe Checkout price at payment time governs.</p>
      </div>
    </main>
  );
}

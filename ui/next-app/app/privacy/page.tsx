export const metadata = { title: "Privacy Policy — Aura Amber" };

export default function PrivacyPage() {
  return (
    <main className="mx-auto max-w-3xl p-6 text-sm leading-6 text-neutral-300">
      <h1 className="mb-4 text-2xl font-semibold text-neutral-100">Privacy Policy</h1>
      <p className="text-neutral-500">Last updated: 10 Sep 2026</p>
      <section className="mt-6 space-y-3">
        <p>Aura Amber migrates store data (products, orders, customers) into <strong>your own PostgreSQL</strong>. We do not sell personal data.</p>
        <h2 className="text-base font-semibold text-neutral-100">Data we process</h2>
        <ul className="list-disc pl-5 text-neutral-400">
          <li>Source store records you explicitly connect (WordPress / Shopify / WooCommerce / Magento).</li>
          <li>Billing contact (email, login) via Stripe; GitHub installation IDs via GitHub App.</li>
          <li>Operational logs (job status, error messages).</li>
        </ul>
        <h2 className="text-base font-semibold text-neutral-100">Processors</h2>
        <p className="text-neutral-400">PostgreSQL (self-hosted/Docker), Redis/BullMQ (job queue), Stripe (billing), GitHub (marketplace), Vercel (hosting).</p>
        <h2 className="text-base font-semibold text-neutral-100">Retention & rights</h2>
        <p className="text-neutral-400">Migration data lives in your database until you delete it. Contact the repo owner to request export/deletion of billing contact data.</p>
      </section>
    </main>
  );
}

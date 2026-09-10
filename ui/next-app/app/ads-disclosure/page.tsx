export const metadata = { title: "Ads Disclosure — Aura Amber" };

export default function AdsDisclosurePage() {
  return (
    <main className="mx-auto max-w-3xl p-6 text-sm leading-6 text-neutral-300">
      <h1 className="mb-4 text-2xl font-semibold text-neutral-100">Ads & Affiliate Disclosure</h1>
      <p className="text-neutral-500">Last updated: 10 Sep 2026</p>
      <div className="mt-6 space-y-3 text-neutral-400">
        <p>Aura Amber itself serves no third-party ads. If marketing pages or docs link to paid services (hosting, Stripe, Vercel), those links may be affiliate or referral links.</p>
        <p>Sponsored content, if any, will be labeled “Sponsored”. Ad partners on your own storefront must be disclosed in your store privacy/cookie pages.</p>
      </div>
    </main>
  );
}

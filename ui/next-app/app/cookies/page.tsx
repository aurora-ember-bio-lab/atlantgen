export const metadata = { title: "Cookie Policy — Aura Amber" };

export default function CookiesPage() {
  return (
    <main className="mx-auto max-w-3xl p-6 text-sm leading-6 text-neutral-300">
      <h1 className="mb-4 text-2xl font-semibold text-neutral-100">Cookie Policy</h1>
      <p className="text-neutral-500">Last updated: 10 Sep 2026</p>
      <div className="mt-6 space-y-3 text-neutral-400">
        <p>We use only strictly-necessary cookies (session, CSRF, consent flag). No advertising trackers are set by Aura Amber itself.</p>
        <p>If you embed analytics or ads on your deployed storefront, disclose them here and gate them behind consent in your own banner.</p>
        <ul className="list-disc pl-5">
          <li><strong>Necessary:</strong> auth session, cookie-consent choice (12 months).</li>
          <li><strong>Optional (only if you enable):</strong> analytics, ad measurement.</li>
        </ul>
      </div>
    </main>
  );
}

function Field({ label, placeholder, value }: { label: string; placeholder?: string; value?: string }) {
  return (
    <label className="block">
      <span className="text-xs font-medium uppercase tracking-wide text-neutral-500">
        {label}
      </span>
      <input
        type="text"
        readOnly
        defaultValue={value}
        placeholder={placeholder}
        className="mt-1 w-full rounded-md border border-neutral-800 bg-neutral-950 px-3 py-2 text-sm text-neutral-200 placeholder:text-neutral-600 focus:border-amber-500 focus:outline-none"
      />
    </label>
  );
}

export default function SettingsPage() {
  return (
    <div className="mx-auto max-w-3xl">
      <div className="mb-8">
        <h1 className="text-2xl font-semibold text-neutral-100">Settings</h1>
        <p className="mt-1 text-sm text-neutral-500">
          Project and deployment configuration.
        </p>
      </div>

      <section className="mb-6 rounded-lg border border-neutral-800 bg-neutral-900/40 p-5">
        <h2 className="mb-4 text-sm font-semibold text-neutral-100">Project</h2>
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <Field label="Project name" value="aura-amber-saas" />
          <Field label="GitHub org" value="cargounetcom" />
          <Field label="Vercel team" value="aurora-ember-cyber" />
          <Field label="Environment" value="production" />
        </div>
      </section>

      <section className="mb-6 rounded-lg border border-neutral-800 bg-neutral-900/40 p-5">
        <h2 className="mb-4 text-sm font-semibold text-neutral-100">API keys</h2>
        <div className="grid grid-cols-1 gap-4">
          <Field label="Database URL" placeholder="postgres://…" />
          <Field label="GitHub App private key" placeholder="Not configured" />
        </div>
        <p className="mt-3 text-xs text-neutral-500">
          Keys are read-only here — set them via environment variables (.env), not this UI.
        </p>
      </section>

      <section className="rounded-lg border border-red-900/40 bg-red-950/10 p-5">
        <h2 className="mb-2 text-sm font-semibold text-red-400">Danger zone</h2>
        <p className="mb-4 text-sm text-neutral-400">
          Remove all migration history and connector credentials.
        </p>
        <button className="rounded-md border border-red-800 px-3 py-2 text-sm font-medium text-red-400 transition-colors hover:bg-red-900/20">
          Reset project
        </button>
      </section>
    </div>
  );
}

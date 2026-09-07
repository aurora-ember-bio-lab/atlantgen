import {
  connectorStateLabels,
  connectorStateStyles,
  mockConnectors,
} from "@/lib/mock-connectors";

function formatDate(iso?: string) {
  if (!iso) return "Never";
  return new Date(iso).toLocaleString("en-IE", {
    day: "2-digit",
    month: "short",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export default function ConnectorsPage() {
  return (
    <div className="mx-auto max-w-6xl">
      <div className="mb-8">
        <h1 className="text-2xl font-semibold text-neutral-100">Connectors</h1>
        <p className="mt-1 text-sm text-neutral-500">
          Source platforms available for migration.
        </p>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        {mockConnectors.map((c) => (
          <div
            key={c.id}
            className="rounded-lg border border-neutral-800 bg-neutral-900/40 p-5"
          >
            <div className="flex items-start justify-between">
              <h2 className="text-base font-semibold text-neutral-100">{c.name}</h2>
              <span
                className={`inline-block rounded-full border px-2 py-0.5 text-xs font-medium ${connectorStateStyles[c.state]}`}
              >
                {connectorStateLabels[c.state]}
              </span>
            </div>
            <p className="mt-2 text-sm text-neutral-400">{c.description}</p>
            <p className="mt-4 text-xs text-neutral-500">
              Last sync: {formatDate(c.lastSync)}
            </p>
            <button className="mt-4 w-full rounded-md border border-neutral-700 px-3 py-2 text-sm font-medium text-neutral-200 transition-colors hover:border-amber-500 hover:text-amber-400">
              {c.state === "connected" ? "Reconfigure" : "Connect"}
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}

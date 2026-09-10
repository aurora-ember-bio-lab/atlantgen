import { Migration, sourceLabels, statusStyles } from "@/lib/migration-types";

function formatDate(iso: string) {
  return new Date(iso).toLocaleString("en-IE", {
    day: "2-digit",
    month: "short",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export default function MigrationTable({ migrations }: { migrations: Migration[] }) {
  return (
    <div className="overflow-hidden rounded-lg border border-neutral-800">
      <table className="w-full text-left text-sm">
        <thead className="bg-neutral-900/60 text-xs uppercase tracking-wide text-neutral-500">
          <tr>
            <th className="px-4 py-3 font-medium">Site</th>
            <th className="px-4 py-3 font-medium">Source</th>
            <th className="px-4 py-3 font-medium">Target</th>
            <th className="px-4 py-3 font-medium">Progress</th>
            <th className="px-4 py-3 font-medium">Status</th>
            <th className="px-4 py-3 font-medium">Started</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-neutral-800">
          {migrations.map((m) => {
            const pct = m.recordsTotal
              ? Math.round((m.recordsMigrated / m.recordsTotal) * 100)
              : 0;
            return (
              <tr key={m.id} className="bg-neutral-950/40 hover:bg-neutral-900/40">
                <td className="px-4 py-3 font-medium text-neutral-100">{m.name}</td>
                <td className="px-4 py-3 text-neutral-400">{sourceLabels[m.source]}</td>
                <td className="px-4 py-3 text-neutral-400">{m.target}</td>
                <td className="px-4 py-3">
                  <div className="flex items-center gap-2">
                    <div className="h-1.5 w-24 overflow-hidden rounded-full bg-neutral-800">
                      <div
                        className="h-full rounded-full bg-amber-500"
                        style={{ width: `${pct}%` }}
                      />
                    </div>
                    <span className="text-xs text-neutral-500">{pct}%</span>
                  </div>
                </td>
                <td className="px-4 py-3">
                  <span
                    className={`inline-block rounded-full border px-2 py-0.5 text-xs font-medium capitalize ${statusStyles[m.status]}`}
                  >
                    {m.status}
                  </span>
                </td>
                <td className="px-4 py-3 text-neutral-500">{formatDate(m.startedAt)}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

"use client";

import { get } from "@/lib/api-client";
import { useEffect, useState } from "react";
import { Migration } from "@/lib/migration-types";

export default function MigrationsPage() {
  const [jobs, setJobs] = useState<Migration[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function loadJobs() {
    try {
      setLoading(true);
      setError(null);
      const data = await get<Migration[]>("/api/migrations");
      setJobs(data);
    } catch (err: any) {
      setError(err.message || "Unknown error");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadJobs();
    const id = setInterval(loadJobs, 2000);
    return () => clearInterval(id);
  }, []);

  const statusColors: Record<string, string> = {
    queued: "bg-neutral-800 text-neutral-300 border-neutral-700",
    running: "bg-amber-500/10 text-amber-400 border-amber-500/30",
    completed: "bg-emerald-500/10 text-emerald-400 border-emerald-500/30",
    failed: "bg-red-500/10 text-red-400 border-red-500/30",
  };

  return (
    <div className="mx-auto max-w-6xl">
      <h1 className="mb-4 text-2xl font-semibold text-neutral-100">
        Migrations
      </h1>

      {loading && <p className="text-neutral-400">Loading…</p>}
      {error && <p className="text-red-400">{error}</p>}

      <div className="overflow-hidden rounded-lg border border-neutral-800">
        <table className="w-full text-left text-sm">
          <thead className="bg-neutral-900/60 text-xs uppercase tracking-wide text-neutral-500">
            <tr>
              <th className="px-4 py-3 font-medium">ID</th>
              <th className="px-4 py-3 font-medium">Source</th>
              <th className="px-4 py-3 font-medium">Status</th>
              <th className="px-4 py-3 font-medium">Progress</th>
              <th className="px-4 py-3 font-medium">Created</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-neutral-800">
            {jobs.map((job) => (
              <tr key={job.id} className="bg-neutral-950/40 hover:bg-neutral-900/40">
                <td className="px-4 py-3 text-neutral-100">{job.id}</td>
                <td className="px-4 py-3 text-neutral-400">{job.source}</td>
                <td className="px-4 py-3">
                  <span
                    className={`inline-block rounded-full border px-2 py-0.5 text-xs font-medium capitalize ${statusColors[job.status] ?? ""}`}
                  >
                    {job.status}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <div className="flex items-center gap-2">
                    <div className="h-1.5 w-24 overflow-hidden rounded-full bg-neutral-800">
                      <div
                        className="h-full rounded-full bg-amber-500"
                        style={{ width: `${(job.recordsTotal ? Math.round((job.recordsMigrated / job.recordsTotal) * 100) : 0)}%` }}
                      />
                    </div>
                    <span className="text-xs text-neutral-500">
                      {job.recordsTotal
                        ? Math.round((job.recordsMigrated / job.recordsTotal) * 100)
                        : 0}
                      %
                    </span>
                  </div>
                </td>
                <td className="px-4 py-3 text-neutral-500">
                  {new Date(job.startedAt).toLocaleString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {jobs.length === 0 && !loading && !error && (
        <p className="mt-4 text-center text-neutral-500">
          No jobs yet. Trigger a migration from your dashboard.
        </p>
      )}
    </div>
  );
}

import StatCard from "@/components/StatCard";
import MigrationTable from "@/components/MigrationTable";
import { mockMigrations, summaryStats } from "@/lib/mock-data";

export default function DashboardPage() {
  return (
    <div className="mx-auto max-w-6xl">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-neutral-100">
            Migration Dashboard
          </h1>
          <p className="mt-1 text-sm text-neutral-500">
            Status: <span className="text-amber-400">Ready</span>
          </p>
        </div>
        <button className="rounded-md bg-amber-500 px-4 py-2 text-sm font-semibold text-neutral-950 transition-colors hover:bg-amber-400">
          + New Migration
        </button>
      </div>

      <div className="mb-8 grid grid-cols-2 gap-4 md:grid-cols-4">
        <StatCard label="Active" value={summaryStats.active} tone="amber" />
        <StatCard label="Completed" value={summaryStats.completed} tone="emerald" />
        <StatCard label="Failed" value={summaryStats.failed} tone="red" />
        <StatCard
          label="Records Migrated"
          value={summaryStats.recordsMigrated.toLocaleString("en-IE")}
        />
      </div>

      <h2 className="mb-3 text-sm font-medium uppercase tracking-wide text-neutral-500">
        Recent Migrations
      </h2>
      <MigrationTable migrations={mockMigrations} />
    </div>
  );
}

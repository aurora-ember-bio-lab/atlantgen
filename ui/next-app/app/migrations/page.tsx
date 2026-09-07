import MigrationTable from "@/components/MigrationTable";
import { mockMigrations } from "@/lib/mock-data";

export default function MigrationsPage() {
  return (
    <div className="mx-auto max-w-6xl">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-neutral-100">Migrations</h1>
          <p className="mt-1 text-sm text-neutral-500">
            All migration jobs across every connected site.
          </p>
        </div>
        <button className="rounded-md bg-amber-500 px-4 py-2 text-sm font-semibold text-neutral-950 transition-colors hover:bg-amber-400">
          + New Migration
        </button>
      </div>

      <MigrationTable migrations={mockMigrations} />
    </div>
  );
}

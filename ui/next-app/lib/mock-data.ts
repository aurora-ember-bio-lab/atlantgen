export type SourceType = "wordpress" | "shopify" | "woocommerce" | "magento";

export type MigrationStatus = "queued" | "running" | "completed" | "failed";

export interface Migration {
  id: string;
  name: string;
  source: SourceType;
  target: string;
  status: MigrationStatus;
  recordsMigrated: number;
  recordsTotal: number;
  startedAt: string;
}

export const sourceLabels: Record<SourceType, string> = {
  wordpress: "WordPress",
  shopify: "Shopify",
  woocommerce: "WooCommerce",
  magento: "Magento",
};

export const statusStyles: Record<MigrationStatus, string> = {
  queued: "bg-neutral-800 text-neutral-300 border-neutral-700",
  running: "bg-amber-500/10 text-amber-400 border-amber-500/30",
  completed: "bg-emerald-500/10 text-emerald-400 border-emerald-500/30",
  failed: "bg-red-500/10 text-red-400 border-red-500/30",
};

export const mockMigrations: Migration[] = [
  {
    id: "mig_1001",
    name: "familywhole.com",
    source: "woocommerce",
    target: "Next.js + Medusa",
    status: "running",
    recordsMigrated: 3420,
    recordsTotal: 5200,
    startedAt: "2026-09-05T14:12:00Z",
  },
  {
    id: "mig_1000",
    name: "ticketwhole.com",
    source: "wordpress",
    target: "Next.js + Medusa",
    status: "completed",
    recordsMigrated: 1180,
    recordsTotal: 1180,
    startedAt: "2026-09-03T09:40:00Z",
  },
  {
    id: "mig_0999",
    name: "formulajewelry.com",
    source: "shopify",
    target: "Next.js + Medusa",
    status: "queued",
    recordsMigrated: 0,
    recordsTotal: 940,
    startedAt: "2026-09-06T08:00:00Z",
  },
  {
    id: "mig_0998",
    name: "legacy-store-import",
    source: "magento",
    target: "Next.js + Medusa",
    status: "failed",
    recordsMigrated: 210,
    recordsTotal: 2600,
    startedAt: "2026-09-02T18:05:00Z",
  },
];

export const summaryStats = {
  active: mockMigrations.filter((m) => m.status === "running").length,
  completed: mockMigrations.filter((m) => m.status === "completed").length,
  failed: mockMigrations.filter((m) => m.status === "failed").length,
  recordsMigrated: mockMigrations.reduce((sum, m) => sum + m.recordsMigrated, 0),
};

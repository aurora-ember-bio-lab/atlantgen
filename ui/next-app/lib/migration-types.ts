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

import type { RawRecord } from "../connectors/types.js";

export interface NormalizedProduct {
  sourceType: string;
  sourceId: string;
  title: string;
  priceCents: number | null;
}

/**
 * Maps a source-specific product shape onto the common Aura Amber schema
 * (see internal/migration/schema.hcl -> table "products").
 */
export function normalizeData(raw: RawRecord[], sourceType: string): NormalizedProduct[] {
  return raw.map((record) => normalizeOne(record, sourceType));
}

function normalizeOne(record: RawRecord, sourceType: string): NormalizedProduct {
  switch (sourceType) {
    case "wordpress":
      return {
        sourceType,
        sourceId: String(record.id),
        title: extractTitle(record.title) ?? "",
        priceCents: null,
      };
    case "shopify":
      return {
        sourceType,
        sourceId: String(record.id),
        title: String(record.title ?? ""),
        priceCents: toCents((record.variants as RawRecord[] | undefined)?.[0]?.price),
      };
    case "woocommerce":
      return {
        sourceType,
        sourceId: String(record.id),
        title: String(record.name ?? ""),
        priceCents: toCents(record.price),
      };
    case "magento":
      return {
        sourceType,
        sourceId: String(record.id),
        title: String(record.name ?? ""),
        priceCents: toCents(record.price),
      };
    default:
      throw new Error(`normalizeData: unknown sourceType "${sourceType}"`);
  }
}

function extractTitle(title: unknown): string | null {
  if (title && typeof title === "object" && "rendered" in title) {
    return String((title as { rendered: unknown }).rendered);
  }
  return typeof title === "string" ? title : null;
}

function toCents(price: unknown): number | null {
  if (price === undefined || price === null || price === "") return null;
  const value = Number(price);
  return Number.isFinite(value) ? Math.round(value * 100) : null;
}

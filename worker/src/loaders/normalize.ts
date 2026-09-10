import type { RawRecord } from "../connectors/types.js";

export interface NormalizedProduct {
  sourceType: string;
  sourceId: string;
  title: string;
  priceCents: number | null;
}

export interface NormalizedOrder {
  sourceType: string;
  sourceId: string;
  totalCents: number;
  createdAt: string;
}

export interface NormalizedCustomer {
  sourceType: string;
  sourceId: string;
  email: string;
}

/**
 * Maps source-specific records onto the common Aura Amber schema
 * (see internal/migration/schema.hcl: products, orders, customers).
 */
export function normalizeData(
  raw: RawRecord[],
  sourceType: string
): { products: NormalizedProduct[]; orders: NormalizedOrder[]; customers: NormalizedCustomer[] } {
  const products: NormalizedProduct[] = [];
  const orders: NormalizedOrder[] = [];
  const customers: NormalizedCustomer[] = [];

  for (const record of raw) {
    switch (sourceType) {
      case "wordpress": {
        products.push({
          sourceType,
          sourceId: String(record.id),
          title: extractTitle(record.title) ?? "",
          priceCents: null,
        });
        if (record.type === "product" || record.type === "shop_order") {
          orders.push({
            sourceType,
            sourceId: String(record.id),
            totalCents: record.total ? Math.round(Number(record.total) * 100) : 0,
            createdAt: (record.date as string) || new Date().toISOString(),
          });
        }
        if (record.type === "user" || record.type === "customer") {
          customers.push({
            sourceType,
            sourceId: String(record.id),
            email: String(record.email ?? record.user_email ?? ""),
          });
        }
        break;
      }
      case "shopify": {
        const variants = record.variants as RawRecord[] | undefined;
        const v0 = variants?.[0] as Record<string, unknown> | undefined;
        products.push({
          sourceType,
          sourceId: String(record.id),
          title: String(record.title ?? ""),
          priceCents: toCents(v0?.price),
        });
        if (variants && variants.length > 0) {
          orders.push({
            sourceType,
            sourceId: String(record.id),
            totalCents: toCents((record.total_price as Record<string, unknown>)?.price as unknown as number) ?? 0,
            createdAt: (record.created_at as string) || new Date().toISOString(),
          });
        }
        const cust = (record.customer as Record<string, unknown>) || {};
        if (cust.id) {
          customers.push({
            sourceType,
            sourceId: String(cust.id),
            email: String((cust.email as string) ?? ""),
          });
        }
        break;
      }
      case "woocommerce": {
        products.push({
          sourceType,
          sourceId: String(record.id),
          title: String(record.name ?? ""),
          priceCents: toCents(record.price),
        });
        if ((record as Record<string, unknown>).parent !== undefined) {
          const rec = record as Record<string, unknown>;
          orders.push({
            sourceType,
            sourceId: String(record.id),
            totalCents: toCents(rec.total as unknown as number) ?? 0,
            createdAt: (rec.date_created as string) || new Date().toISOString(),
          });
        }
        if ((record as Record<string, unknown>).email) {
          const rec = record as Record<string, unknown>;
          customers.push({
            sourceType,
            sourceId: String(record.id),
            email: String(rec.email as string),
          });
        }
        break;
      }
      case "magento": {
        products.push({
          sourceType,
          sourceId: String(record.id),
          title: String(record.name ?? ""),
          priceCents: toCents(record.price),
        });
        if ((record as Record<string, unknown>).status) {
          const rec = record as Record<string, unknown>;
          orders.push({
            sourceType,
            sourceId: String(record.id),
            totalCents: toCents(rec.grand_total as unknown as number) ?? 0,
            createdAt: (rec.created_at as string) || new Date().toISOString(),
          });
        }
        if ((record as Record<string, unknown>).email) {
          const rec = record as Record<string, unknown>;
          customers.push({
            sourceType,
            sourceId: String(record.id),
            email: String(rec.email as string),
          });
        }
        break;
      }
      default:
        throw new Error(`normalizeData: unknown sourceType "${sourceType}"`);
    }
  }

  return { products, orders, customers };
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

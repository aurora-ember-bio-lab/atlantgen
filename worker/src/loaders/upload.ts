import { Client } from "pg";
import type { NormalizedProduct, NormalizedOrder, NormalizedCustomer } from "./normalize.js";

/**
 * Upserts normalized data into targetDbUrl's tables (products, orders, customers).
 * Upsert key is (source_type, source_id).
 */
export async function uploadData(
  products: NormalizedProduct[],
  orders: NormalizedOrder[],
  customers: NormalizedCustomer[],
  targetDbUrl: string
): Promise<void> {
  const client = new Client({ connectionString: targetDbUrl });
  await client.connect();

  try {
    await client.query("BEGIN");

    if (products.length > 0) {
      for (const record of products) {
        await client.query(
          `INSERT INTO products (id, title, source_type, source_id, price_cents)
           VALUES (gen_random_uuid(), $1, $2, $3, $4)
           ON CONFLICT (source_type, source_id)
           DO UPDATE SET title = EXCLUDED.title, price_cents = EXCLUDED.price_cents`,
          [record.title, record.sourceType, record.sourceId, record.priceCents]
        );
      }
    }

    if (orders.length > 0) {
      for (const record of orders) {
        await client.query(
          `INSERT INTO orders (id, source_type, source_id, total_cents, created_at)
           VALUES (gen_random_uuid(), $1, $2, $3, $4)
           ON CONFLICT (source_type, source_id)
           DO UPDATE SET total_cents = EXCLUDED.total_cents`,
          [record.sourceType, record.sourceId, record.totalCents, record.createdAt]
        );
      }
    }

    if (customers.length > 0) {
      for (const record of customers) {
        await client.query(
          `INSERT INTO customers (id, email, source_type, source_id)
           VALUES (gen_random_uuid(), $1, $2, $3)
           ON CONFLICT (source_type, source_id)
           DO UPDATE SET email = EXCLUDED.email`,
          [record.email, record.sourceType, record.sourceId]
        );
      }
    }

    await client.query("COMMIT");
  } catch (err) {
    await client.query("ROLLBACK");
    throw err;
  } finally {
    await client.end();
  }
}

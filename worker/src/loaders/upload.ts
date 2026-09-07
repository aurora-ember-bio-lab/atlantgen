import { Client } from "pg";
import type { NormalizedProduct } from "./normalize.js";

/**
 * Upserts normalized products into targetDbUrl's `products` table (schema
 * defined in internal/migration/schema.hcl, kept in sync by the Atlas
 * Schema Orchestrator). Upsert key is (source_type, source_id).
 */
export async function uploadData(records: NormalizedProduct[], targetDbUrl: string): Promise<void> {
  if (records.length === 0) return;

  const client = new Client({ connectionString: targetDbUrl });
  await client.connect();

  try {
    await client.query("BEGIN");
    for (const record of records) {
      await client.query(
        `INSERT INTO products (id, title, source_type, source_id, price_cents)
         VALUES (gen_random_uuid(), $1, $2, $3, $4)
         ON CONFLICT (source_type, source_id)
         DO UPDATE SET title = EXCLUDED.title, price_cents = EXCLUDED.price_cents`,
        [record.title, record.sourceType, record.sourceId, record.priceCents],
      );
    }
    await client.query("COMMIT");
  } catch (err) {
    await client.query("ROLLBACK");
    throw err;
  } finally {
    await client.end();
  }
}

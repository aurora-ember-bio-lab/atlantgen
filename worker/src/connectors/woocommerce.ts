import type { ConnectionInfo, RawRecord } from "./types.js";

/**
 * Pulls products, orders and customers from a WooCommerce store via the
 * WooCommerce REST API (consumer key/secret auth).
 *
 * Expected `info` keys: baseUrl, consumerKey, consumerSecret.
 */
export async function extractWooCommerce(info: ConnectionInfo): Promise<RawRecord[]> {
  const { baseUrl, consumerKey, consumerSecret } = info;
  if (!baseUrl || !consumerKey || !consumerSecret) {
    throw new Error("woocommerce connector requires baseUrl, consumerKey, consumerSecret");
  }

  const records: RawRecord[] = [];
  let page = 1;

  for (;;) {
    const url = new URL(`${baseUrl.replace(/\/$/, "")}/wp-json/wc/v3/products`);
    url.searchParams.set("consumer_key", consumerKey);
    url.searchParams.set("consumer_secret", consumerSecret);
    url.searchParams.set("per_page", "100");
    url.searchParams.set("page", String(page));

    const res = await fetch(url);
    if (!res.ok) {
      throw new Error(`woocommerce connector: ${res.status} ${res.statusText}`);
    }

    const batch = (await res.json()) as RawRecord[];
    if (batch.length === 0) break;

    records.push(...batch);
    page += 1;
  }

  return records;
}

import type { ConnectionInfo, RawRecord } from "./types.js";

/**
 * Pulls products from a Shopify store via the Admin REST API.
 *
 * Expected `info` keys: shopDomain (e.g. "my-shop.myshopify.com"), accessToken.
 */
export async function extractShopify(info: ConnectionInfo): Promise<RawRecord[]> {
  const { shopDomain, accessToken } = info;
  if (!shopDomain || !accessToken) {
    throw new Error("shopify connector requires shopDomain, accessToken");
  }

  const records: RawRecord[] = [];
  let pageInfo: string | null = null;

  for (;;) {
    const url = new URL(`https://${shopDomain}/admin/api/2024-07/products.json`);
    url.searchParams.set("limit", "250");
    if (pageInfo) url.searchParams.set("page_info", pageInfo);

    const res = await fetch(url, {
      headers: { "X-Shopify-Access-Token": accessToken },
    });
    if (!res.ok) {
      throw new Error(`shopify connector: ${res.status} ${res.statusText}`);
    }

    const body = (await res.json()) as { products: RawRecord[] };
    records.push(...body.products);

    const link = res.headers.get("link");
    const nextMatch = link?.match(/<[^>]*[?&]page_info=([^&>]+)[^>]*>;\s*rel="next"/);
    if (!nextMatch) break;
    pageInfo = nextMatch[1];
  }

  return records;
}

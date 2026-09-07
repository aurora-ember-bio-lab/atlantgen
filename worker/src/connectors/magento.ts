import type { ConnectionInfo, RawRecord } from "./types.js";

/**
 * Pulls products from a Magento 2 store via the REST API.
 *
 * Expected `info` keys: baseUrl, accessToken (integration/admin token).
 */
export async function extractMagento(info: ConnectionInfo): Promise<RawRecord[]> {
  const { baseUrl, accessToken } = info;
  if (!baseUrl || !accessToken) {
    throw new Error("magento connector requires baseUrl, accessToken");
  }

  const records: RawRecord[] = [];
  let currentPage = 1;
  const pageSize = 100;

  for (;;) {
    const url = new URL(`${baseUrl.replace(/\/$/, "")}/rest/V1/products`);
    url.searchParams.set("searchCriteria[pageSize]", String(pageSize));
    url.searchParams.set("searchCriteria[currentPage]", String(currentPage));

    const res = await fetch(url, {
      headers: { Authorization: `Bearer ${accessToken}` },
    });
    if (!res.ok) {
      throw new Error(`magento connector: ${res.status} ${res.statusText}`);
    }

    const body = (await res.json()) as { items: RawRecord[] };
    if (body.items.length === 0) break;

    records.push(...body.items);
    if (body.items.length < pageSize) break;
    currentPage += 1;
  }

  return records;
}

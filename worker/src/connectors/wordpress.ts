import type { ConnectionInfo, RawRecord } from "./types.js";

/**
 * Pulls posts/pages/products from a WordPress site via the REST API
 * (wp-json/wp/v2/...), authenticated with an application password.
 *
 * Expected `info` keys: baseUrl, username, appPassword.
 */
export async function extractWordPress(info: ConnectionInfo): Promise<RawRecord[]> {
  const { baseUrl, username, appPassword } = info;
  if (!baseUrl || !username || !appPassword) {
    throw new Error("wordpress connector requires baseUrl, username, appPassword");
  }

  const auth = Buffer.from(`${username}:${appPassword}`).toString("base64");
  const records: RawRecord[] = [];

  let page = 1;
  for (;;) {
    const res = await fetch(`${baseUrl.replace(/\/$/, "")}/wp-json/wp/v2/posts?per_page=100&page=${page}`, {
      headers: { Authorization: `Basic ${auth}` },
    });

    if (res.status === 400) break; // WordPress returns 400 past the last page
    if (!res.ok) {
      throw new Error(`wordpress connector: ${res.status} ${res.statusText}`);
    }

    const batch = (await res.json()) as RawRecord[];
    if (batch.length === 0) break;

    records.push(...batch);
    page += 1;
  }

  return records;
}

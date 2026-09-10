import { Worker, type Job } from "bullmq";
import { connection, type MigrationJobData, type SourceType } from "./queue.js";
import { startServer } from "./server.js";
import { extractWordPress } from "./connectors/wordpress.js";
import { extractShopify } from "./connectors/shopify.js";
import { extractWooCommerce } from "./connectors/woocommerce.js";
import { extractMagento } from "./connectors/magento.js";
import { normalizeData } from "./loaders/normalize.js";
import { uploadData } from "./loaders/upload.js";
import type { ConnectionInfo, RawRecord } from "./connectors/types.js";

const extractors: Record<SourceType, (info: ConnectionInfo) => Promise<RawRecord[]>> = {
  wordpress: extractWordPress,
  shopify: extractShopify,
  woocommerce: extractWooCommerce,
  magento: extractMagento,
};

const worker = new Worker<MigrationJobData>(
  "migration",
  async (job: Job<MigrationJobData>) => {
    const { sourceType, sourceConnectionInfo, targetDbUrl } = job.data;

    const extract = extractors[sourceType];
    if (!extract) {
      throw new Error(`Unknown sourceType: ${sourceType}`);
    }

    await job.updateProgress(10);
    const raw = await extract(sourceConnectionInfo);

    await job.updateProgress(50);
    const { products, orders, customers } = normalizeData(raw, sourceType);

    await job.updateProgress(75);
    await uploadData(products, orders, customers, targetDbUrl);

    await job.updateProgress(100);
    return { status: "completed", records: { products: products.length, orders: orders.length, customers: customers.length } };
  },
  { connection },
);

worker.on("completed", (job) => {
  console.log(`[migration] job ${job.id} completed`);
});

worker.on("failed", (job, err) => {
  console.error(`[migration] job ${job?.id} failed:`, err.message);
});

startServer();

console.log("Aura Amber migration worker started, waiting for jobs...");

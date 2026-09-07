import { Queue } from "bullmq";
import type { ConnectionInfo } from "./connectors/types.js";

export const connection = {
  host: process.env.REDIS_HOST ?? "localhost",
  port: Number(process.env.REDIS_PORT ?? 6379),
};

export type SourceType = "wordpress" | "shopify" | "woocommerce" | "magento";

export interface MigrationJobData {
  sourceType: SourceType;
  sourceConnectionInfo: ConnectionInfo;
  targetDbUrl: string;
}

export const migrationQueue = new Queue<MigrationJobData>("migration", { connection });

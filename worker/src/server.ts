import { createServer } from "node:http";
import { migrationQueue, type MigrationJobData } from "./queue.js";

const PORT = Number(process.env.WORKER_HTTP_PORT ?? 8090);

/**
 * Small HTTP shim the Go API calls to enqueue a migration job. BullMQ's
 * Redis wire format is Node-specific, so instead of Go writing to Redis
 * directly, it POSTs the job here and this process (which already speaks
 * BullMQ natively) puts it on the real queue.
 */
export function startServer() {
  const server = createServer((req, res) => {
    if (req.method === "GET" && req.url === "/health") {
      res.writeHead(200, { "Content-Type": "application/json" });
      res.end(JSON.stringify({ status: "ok" }));
      return;
    }

    if (req.method === "POST" && req.url === "/enqueue") {
      let body = "";
      req.on("data", (chunk) => (body += chunk));
      req.on("end", async () => {
        try {
          const payload = JSON.parse(body) as { name: string } & MigrationJobData;
          if (!payload.name || !payload.sourceType || !payload.targetDbUrl) {
            res.writeHead(400, { "Content-Type": "application/json" });
            res.end(JSON.stringify({ error: "name, sourceType and targetDbUrl are required" }));
            return;
          }

          const job = await migrationQueue.add(payload.name, {
            sourceType: payload.sourceType,
            sourceConnectionInfo: payload.sourceConnectionInfo ?? {},
            targetDbUrl: payload.targetDbUrl,
          });

          res.writeHead(201, { "Content-Type": "application/json" });
          res.end(JSON.stringify({ jobId: job.id }));
        } catch (err) {
          res.writeHead(400, { "Content-Type": "application/json" });
          res.end(JSON.stringify({ error: (err as Error).message }));
        }
      });
      return;
    }

    if (req.method === "GET" && req.url?.startsWith("/jobs/")) {
      const id = req.url.slice("/jobs/".length);
      migrationQueue
        .getJob(id)
        .then(async (job) => {
          if (!job) {
            res.writeHead(404, { "Content-Type": "application/json" });
            res.end(JSON.stringify({ error: "not found" }));
            return;
          }
          const state = await job.getState();
          res.writeHead(200, { "Content-Type": "application/json" });
          res.end(JSON.stringify({ id: job.id, state, progress: job.progress }));
        })
        .catch((err) => {
          res.writeHead(500, { "Content-Type": "application/json" });
          res.end(JSON.stringify({ error: (err as Error).message }));
        });
      return;
    }

    res.writeHead(404, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ error: "not found" }));
  });

  server.listen(PORT, () => {
    console.log(`Aura Amber worker HTTP listening on :${PORT}`);
  });

  return server;
}

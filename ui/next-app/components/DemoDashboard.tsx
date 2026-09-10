"use client";

import { useState, useEffect } from "react";

const mockMigrations = [
  { id: "mig_001", name: "wp-store-migration", source: "wordpress", status: "completed", progress: 100, records: 12847, startedAt: "2026-09-10T14:23:00Z" },
  { id: "mig_002", name: "woo-sync-q3", source: "woocommerce", status: "running", progress: 67, records: 8234, startedAt: "2026-09-10T15:45:00Z" },
  { id: "mig_003", name: "shopify-cutover", source: "shopify", status: "queued", progress: 0, records: 0, startedAt: "2026-09-10T16:10:00Z" },
  { id: "mig_004", name: "magento-replatform", source: "magento", status: "running", progress: 34, records: 3102, startedAt: "2026-09-10T16:02:00Z" },
];

const mockConnectors = [
  { id: "wordpress", name: "WordPress", state: "connected", lastSync: "2 min ago" },
  { id: "woocommerce", name: "WooCommerce", state: "connected", lastSync: "5 min ago" },
  { id: "shopify", name: "Shopify", state: "not_connected", lastSync: "Never" },
  { id: "magento", name: "Magento", state: "error", lastSync: "3 days ago" },
];

const statusColors: Record<string, string> = {
  queued: "bg-neutral-700 text-neutral-300",
  running: "bg-amber-500/15 text-amber-400 border border-amber-500/30",
  completed: "bg-emerald-500/15 text-emerald-400 border border-emerald-500/30",
  failed: "bg-red-500/15 text-red-400 border border-red-500/30",
};

const connectorColors: Record<string, string> = {
  connected: "bg-emerald-500/15 text-emerald-400 border border-emerald-500/30",
  not_connected: "bg-neutral-700 text-neutral-400 border border-neutral-600",
  error: "bg-red-500/15 text-red-400 border border-red-500/30",
};

export default function DemoDashboard() {
  const [tick, setTick] = useState(0);
  const [migrations, setMigrations] = useState(mockMigrations);

  useEffect(() => {
    const interval = setInterval(() => {
      setTick((t) => t + 1);
      setMigrations((prev) =>
        prev.map((m) => {
          if (m.status === "running") {
            const newProgress = Math.min(m.progress + Math.random() * 8, 100);
            return {
              ...m,
              progress: Math.round(newProgress),
              records: m.records + Math.floor(Math.random() * 120),
              status: newProgress >= 100 ? "completed" : "running",
            };
          }
          if (m.status === "queued" && tick > 6) {
            return { ...m, status: "running", progress: 5, records: 12 };
          }
          return m;
        })
      );
    }, 1200);
    return () => clearInterval(interval);
  }, [tick]);

  const stats = {
    active: migrations.filter((m) => m.status === "running").length,
    completed: migrations.filter((m) => m.status === "completed").length,
    failed: migrations.filter((m) => m.status === "failed").length,
    records: migrations.reduce((s, m) => s + m.records, 0),
  };

  return (
    <div className="relative rounded-2xl border border-neutral-800/60 bg-neutral-900/40 backdrop-blur-xl shadow-2xl shadow-black/30 overflow-hidden">
      {/* Title bar */}
      <div className="flex items-center gap-2 border-b border-neutral-800/60 px-5 py-3">
        <div className="h-3 w-3 rounded-full bg-red-500/60" />
        <div className="h-3 w-3 rounded-full bg-yellow-500/60" />
        <div className="h-3 w-3 rounded-full bg-green-500/60" />
        <span className="ml-3 text-xs font-medium text-neutral-500">Migration Dashboard</span>
        <span className="ml-auto flex items-center gap-1.5 text-xs text-emerald-400">
          <span className="h-1.5 w-1.5 rounded-full bg-emerald-400 animate-pulse" />
          Live
        </span>
      </div>

      <div className="p-6">
        {/* Stats row */}
        <div className="mb-6 grid grid-cols-4 gap-3">
          {[
            { label: "Active", value: stats.active, color: "text-amber-400" },
            { label: "Completed", value: stats.completed, color: "text-emerald-400" },
            { label: "Failed", value: stats.failed, color: "text-red-400" },
            { label: "Records", value: stats.records.toLocaleString(), color: "text-neutral-100" },
          ].map((s) => (
            <div key={s.label} className="rounded-lg border border-neutral-800/60 bg-neutral-950/50 p-3 text-center">
              <div className={`text-xl font-bold ${s.color}`}>{s.value}</div>
              <div className="text-[10px] uppercase tracking-wider text-neutral-500">{s.label}</div>
            </div>
          ))}
        </div>

        {/* Migrations table */}
        <div className="mb-6 overflow-hidden rounded-lg border border-neutral-800/60">
          <table className="w-full text-left text-xs">
            <thead className="bg-neutral-900/60 text-[10px] uppercase tracking-wider text-neutral-500">
              <tr>
                <th className="px-3 py-2 font-medium">Name</th>
                <th className="px-3 py-2 font-medium">Source</th>
                <th className="px-3 py-2 font-medium">Status</th>
                <th className="px-3 py-2 font-medium">Progress</th>
                <th className="px-3 py-2 font-medium">Records</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-neutral-800/40">
              {migrations.map((m) => (
                <tr key={m.id} className="bg-neutral-950/30 hover:bg-neutral-900/30 transition-colors">
                  <td className="px-3 py-2.5 font-medium text-neutral-200">{m.name}</td>
                  <td className="px-3 py-2.5 text-neutral-400 capitalize">{m.source}</td>
                  <td className="px-3 py-2.5">
                    <span className={`inline-block rounded-full px-2 py-0.5 text-[10px] font-medium capitalize ${statusColors[m.status]}`}>
                      {m.status}
                    </span>
                  </td>
                  <td className="px-3 py-2.5">
                    <div className="flex items-center gap-2">
                      <div className="h-1.5 w-16 overflow-hidden rounded-full bg-neutral-800">
                        <div
                          className="h-full rounded-full bg-amber-500 transition-all duration-500 ease-out"
                          style={{ width: `${m.progress}%` }}
                        />
                      </div>
                      <span className="text-[10px] text-neutral-500 w-7">{m.progress}%</span>
                    </div>
                  </td>
                  <td className="px-3 py-2.5 text-neutral-400">{m.records.toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Connectors row */}
        <div className="grid grid-cols-4 gap-3">
          {mockConnectors.map((c) => (
            <div key={c.id} className="rounded-lg border border-neutral-800/60 bg-neutral-950/50 p-3">
              <div className="flex items-center justify-between">
                <span className="text-xs font-medium text-neutral-200">{c.name}</span>
                <span className={`rounded-full px-1.5 py-0.5 text-[9px] font-medium ${connectorColors[c.state]}`}>
                  {c.state === "connected" ? "● Live" : c.state === "error" ? "● Error" : "○ Off"}
                </span>
              </div>
              <div className="mt-1 text-[10px] text-neutral-500">Last sync: {c.lastSync}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

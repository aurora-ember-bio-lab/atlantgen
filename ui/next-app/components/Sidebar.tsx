"use client";

import { usePathname } from "next/navigation";

const navItems = [
  { label: "Dashboard", href: "/" },
  { label: "Migrations", href: "/migrations" },
  { label: "Connectors", href: "/connectors" },
  { label: "Settings", href: "/settings" },
];

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="hidden w-64 shrink-0 border-r border-neutral-800 bg-neutral-950 px-6 py-8 md:block">
      <div className="mb-10 flex items-center gap-2">
        <div className="flex h-8 w-8 items-center justify-center rounded-md bg-amber-500 text-sm font-bold text-neutral-950">
          A
        </div>
        <span className="text-lg font-semibold tracking-tight text-neutral-100">
          Aura Amber
        </span>
      </div>

      <nav className="space-y-1">
        {navItems.map((item) => {
          const active = pathname === item.href;
          return (
            <a
              key={item.href}
              href={item.href}
              className={`block rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                active
                  ? "bg-amber-500/10 text-amber-400"
                  : "text-neutral-400 hover:bg-neutral-900 hover:text-neutral-100"
              }`}
            >
              {item.label}
            </a>
          );
        })}
      </nav>

      <div className="mt-10 rounded-lg border border-neutral-800 bg-neutral-900/60 p-4">
        <p className="text-xs font-medium text-neutral-400">
          WordPress / Shopify / WooCommerce / Magento
        </p>
        <p className="mt-1 text-xs text-neutral-500">
          → Next.js + MedusaJS + PostgreSQL
        </p>
      </div>

      <div className="mt-6 flex flex-wrap gap-x-3 gap-y-1 text-xs text-neutral-500">
        <a href="/privacy" className="hover:text-amber-400">Privacy</a>
        <a href="/cookies" className="hover:text-amber-400">Cookies</a>
        <a href="/disclaimer" className="hover:text-amber-400">Disclaimer</a>
        <a href="/ads-disclosure" className="hover:text-amber-400">Ads</a>
      </div>
    </aside>
  );
}

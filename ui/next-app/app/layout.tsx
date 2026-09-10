import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Aura Amber — Migrate your store without losing a beat",
  description:
    "WordPress / Shopify / WooCommerce / Magento → Next.js + MedusaJS + PostgreSQL migration platform",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-neutral-950 text-neutral-100 antialiased">
        {children}
      </body>
    </html>
  );
}

import type { Metadata } from "next";
import "./globals.css";
import Sidebar from "@/components/Sidebar";

export const metadata: Metadata = {
  title: "Aura Amber — Migration Dashboard",
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
        <div className="flex min-h-screen">
          <Sidebar />
          <main className="flex-1 px-6 py-8 md:px-10">{children}</main>
        </div>
      </body>
    </html>
  );
}

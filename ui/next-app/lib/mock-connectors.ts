export type ConnectorState = "connected" | "not_connected" | "error";

export interface Connector {
  id: string;
  name: string;
  description: string;
  state: ConnectorState;
  lastSync?: string;
}

export const connectorStateStyles: Record<ConnectorState, string> = {
  connected: "bg-emerald-500/10 text-emerald-400 border-emerald-500/30",
  not_connected: "bg-neutral-800 text-neutral-400 border-neutral-700",
  error: "bg-red-500/10 text-red-400 border-red-500/30",
};

export const connectorStateLabels: Record<ConnectorState, string> = {
  connected: "Connected",
  not_connected: "Not connected",
  error: "Error",
};

export const mockConnectors: Connector[] = [
  {
    id: "wordpress",
    name: "WordPress",
    description: "REST API + WP-CLI export for posts, pages, media and users.",
    state: "connected",
    lastSync: "2026-09-05T14:12:00Z",
  },
  {
    id: "woocommerce",
    name: "WooCommerce",
    description: "Products, orders, customers via the WooCommerce REST API.",
    state: "connected",
    lastSync: "2026-09-05T14:12:00Z",
  },
  {
    id: "shopify",
    name: "Shopify",
    description: "Admin API export for products, collections and orders.",
    state: "not_connected",
  },
  {
    id: "magento",
    name: "Magento",
    description: "REST/GraphQL export for catalog, CMS and customer data.",
    state: "error",
    lastSync: "2026-09-02T18:05:00Z",
  },
];

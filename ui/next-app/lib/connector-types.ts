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

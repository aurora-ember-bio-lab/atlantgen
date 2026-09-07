export interface ConnectionInfo {
  [key: string]: string;
}

/** A single extracted record, shape depends on the source platform. */
export interface RawRecord {
  [key: string]: unknown;
}

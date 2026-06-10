// Mirror of server/api/requestSpec.schema.json (the single source of truth).
// Keep in sync with internal/model/requestSpec.go. The Go-side contract
// round-trip test catches drift.

export type Method = "GET" | "POST" | "PUT" | "PATCH" | "DELETE" | "HEAD" | "OPTIONS";

export interface Param {
  name: string;
  value?: string;
  values?: string[];
  required?: boolean;
  description?: string;
}

export interface Header {
  name: string;
  value: string;
  required?: boolean;
}

export type AuthType = "none" | "bearer" | "basic" | "apiKey";

export interface AuthSpec {
  type: AuthType;
  in?: "header" | "query";
  name?: string;
  prefix?: string;
  valueRef?: string;
  username?: string;
  password?: string;
}

export type Encoding = "none" | "json" | "form" | "multipart" | "raw";

export interface BodySpec {
  encoding: Encoding;
  json?: unknown;
  form?: Param[];
  raw?: string;
  contentType?: string;
}

export interface RequestSpec {
  method: Method;
  baseUrl: string;
  path: string;
  pathParams?: Param[];
  query?: Param[];
  headers?: Header[];
  auth: AuthSpec;
  body: BodySpec;
  confidence: number;
  notes?: string;
}

// --- Wire DTOs (matches server/internal/dto) ---

export interface ExecuteRequest {
  requestSpec: RequestSpec;
  secretRefs?: Record<string, string>;
}

export interface RequestPreview {
  method: string;
  url: string;
  headers: Record<string, string>;
  body?: string;
}

export interface ExecuteResponse {
  status: number;
  durationMs: number;
  responseHeaders: Record<string, string>;
  body: string;
  truncated: boolean;
  requestPreview: RequestPreview;
  consentRequired?: { secret: string; host: string };
}

export interface PairRequest {
  token: string;
}

export interface PairResponse {
  ok: boolean;
  boundOrigin: string;
}

export interface HealthResponse {
  status: "ok";
  version: string;
  paired: boolean;
}

export interface ErrorEnvelope {
  error: { code: string; message: string; details?: unknown };
}

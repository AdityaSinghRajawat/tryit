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

// --- Parse cascade ---

export type ParseSource = "cache" | "spec" | "extractor" | "profile" | "ai";

export interface ParseRequest {
  pageUrl: string;
  scopedMarkdown: string;
  authSectionMarkdown?: string;
  framework?: string;
  structuredHint?: unknown;
  force?: boolean;
}

export interface ParseResponse {
  requestSpec: RequestSpec;
  source: ParseSource;
  confidence: number;
  needsConfirmation: boolean;
}

// --- Secrets vault ---

export type SecretType = "bearer" | "basic" | "apiKey";

export interface SecretInfo {
  name: string;
  type: SecretType;
  hostHint?: string;
}

export interface SecretsListResponse {
  secrets: SecretInfo[];
}

export interface SecretCreateRequest {
  name: string;
  type: SecretType;
  hostHint?: string;
  value?: string;
  username?: string;
  password?: string;
}

export interface SecretCreateResponse {
  name: string;
}

export interface SecretDeleteResponse {
  name: string;
  deleted: boolean;
}

// --- Consent ---

export interface ConsentRequest {
  secret: string;
  host: string;
}

export interface ConsentResponse {
  granted: boolean;
}

// --- Codegen ---

export type GenerateLanguage = "curl" | "python" | "javascript";

export interface GenerateRequest {
  requestSpec: RequestSpec;
  language: GenerateLanguage;
  idiomatic?: boolean;
}

export interface GenerateResponse {
  language: GenerateLanguage;
  code: string;
}

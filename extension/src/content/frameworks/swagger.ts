// Swagger UI detection + extraction. Phase 1's only target framework.
//
// We rely on the standard Swagger UI 3.x / 5.x DOM shape:
//   .swagger-ui .opblock                — one per endpoint
//     .opblock-summary-method           — text: GET / POST / ...
//     .opblock-summary-path             — text: /pet/{petId}
//   .servers select                     — server dropdown with base URLs
//
// We build a RequestSpec hint client-side (Phase 1: no server cascade); the
// panel sends it straight to /execute. The user can edit baseUrl/auth in the
// panel before sending.

import type { RequestSpec, Method } from "../../shared/types";

export function isSwaggerUI(): boolean {
  return !!document.querySelector(".swagger-ui .opblock");
}

export function findEndpointBlocks(): HTMLElement[] {
  return Array.from(document.querySelectorAll<HTMLElement>(".swagger-ui .opblock"));
}

export function findInjectionAnchor(opblock: HTMLElement): HTMLElement | null {
  // The summary row is the cleanest place — sits next to method/path/description.
  return (
    opblock.querySelector<HTMLElement>(".opblock-summary-control") ||
    opblock.querySelector<HTMLElement>(".opblock-summary") ||
    opblock
  );
}

const METHODS: ReadonlySet<Method> = new Set([
  "GET",
  "POST",
  "PUT",
  "PATCH",
  "DELETE",
  "HEAD",
  "OPTIONS",
]);

export function buildRequestSpec(opblock: HTMLElement): RequestSpec | null {
  const methodText = opblock
    .querySelector(".opblock-summary-method")
    ?.textContent?.trim()
    .toUpperCase();
  const pathText = opblock.querySelector(".opblock-summary-path")?.textContent?.trim();
  if (!methodText || !pathText) return null;
  if (!METHODS.has(methodText as Method)) return null;

  const baseUrl = detectBaseUrl();
  if (!baseUrl) return null;

  return {
    method: methodText as Method,
    baseUrl,
    path: pathText.startsWith("/") ? pathText : "/" + pathText,
    pathParams: extractPathParams(pathText),
    query: [],
    headers: [],
    auth: { type: "none" },
    body: { encoding: "none" },
    confidence: 0.95,
    notes: "swagger-ui (Phase 1 client-side hint)",
  };
}

function detectBaseUrl(): string | null {
  // 1. Swagger UI 3/5: <select> inside the servers section.
  const sel = document.querySelector<HTMLSelectElement>(".servers select, select[aria-label='Servers']");
  if (sel?.value) return trimSlash(sel.value);
  // 2. Some Petstore-style pages list the URL as text under .info pre.
  const pre = document.querySelector<HTMLElement>(".info pre");
  if (pre?.textContent) return trimSlash(pre.textContent.trim());
  // 3. Fallback: page origin. The user can override in the panel.
  return trimSlash(window.location.origin);
}

function trimSlash(s: string): string {
  return s.replace(/\/+$/, "");
}

function extractPathParams(path: string): { name: string; required: boolean }[] {
  const out: { name: string; required: boolean }[] = [];
  const re = /\{([^}]+)\}/g;
  let m: RegExpExecArray | null;
  while ((m = re.exec(path)) !== null) out.push({ name: m[1], required: true });
  return out;
}

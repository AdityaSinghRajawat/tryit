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

import type { BodySpec, Param, RequestSpec, Method } from "../../shared/types";

export function isSwaggerUI(): boolean {
  return !!document.querySelector(".swagger-ui, #swagger-ui");
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

  const path = pathText.startsWith("/") ? pathText : "/" + pathText;
  const scraped = scrapeParameters(opblock);
  // Merge path params from the URL pattern (always present) with scraped (may
  // add `description` / example values). The URL pattern is authoritative for
  // existence; scraped enriches.
  const pathParams = mergePathParams(extractPathParams(path), scraped.path);

  return {
    method: methodText as Method,
    baseUrl,
    path,
    pathParams,
    query: scraped.query,
    headers: scraped.header.map((p) => ({
      name: p.name,
      value: p.value ?? "",
      required: p.required,
    })),
    auth: { type: "none" },
    body: inferBody(methodText, scraped.body, scraped.formData),
    confidence: scraped.expanded ? 0.95 : 0.7,
    notes: scraped.expanded
      ? "swagger-ui (Phase 1 client-side hint, parameters scraped)"
      : "swagger-ui (Phase 1 client-side hint, opblock not expanded — params editable)",
  };
}

interface ScrapedParams {
  expanded: boolean;
  path: Param[];
  query: Param[];
  header: Param[];
  body: unknown | null;           // OAS3 application/json schema example, if any
  formData: Param[];              // OAS2 formData params (and OAS3 form)
}

// scrapeParameters reads the (already-rendered) Parameters table inside an
// opblock. Swagger UI lazy-renders these; if the block is collapsed the table
// is absent — we return empty arrays and the user fills the editor manually.
function scrapeParameters(opblock: HTMLElement): ScrapedParams {
  const result: ScrapedParams = { expanded: false, path: [], query: [], header: [], body: null, formData: [] };
  const rows = opblock.querySelectorAll<HTMLElement>(".parameters tbody tr.parameters, .parameters tbody tr");
  if (rows.length === 0) return result;
  result.expanded = true;

  for (const row of Array.from(rows)) {
    const nameEl = row.querySelector(".parameter__name");
    if (!nameEl) continue;
    const rawName = nameEl.textContent?.replace(/\*$/, "").trim();
    if (!rawName) continue;
    const required = nameEl.classList.contains("required") || /\*/.test(nameEl.textContent ?? "");
    const inText = row.querySelector(".parameter__in")?.textContent?.toLowerCase() ?? "";
    const description = row.querySelector(".parameters-col_description p, .parameters-col_description .markdown p")?.textContent?.trim() ?? "";

    // Example value from the small inline input/textarea Swagger renders.
    const example = (
      row.querySelector<HTMLInputElement>(".parameters-col_description input") ??
      row.querySelector<HTMLTextAreaElement>(".parameters-col_description textarea") ??
      row.querySelector<HTMLSelectElement>(".parameters-col_description select")
    )?.value ?? "";

    const p: Param = { name: rawName, required, value: example || undefined, description: description || undefined };

    if (inText.includes("path")) result.path.push(p);
    else if (inText.includes("query")) result.query.push(p);
    else if (inText.includes("header")) result.header.push(p);
    else if (inText.includes("formdata") || inText.includes("form data")) result.formData.push(p);
    else if (inText.includes("body")) {
      // OAS2 single body parameter; example may be a JSON blob.
      if (example.trim().startsWith("{") || example.trim().startsWith("[")) {
        try { result.body = JSON.parse(example); }
        catch { result.body = example; }
      }
    }
  }

  // OAS3: the request body lives under .opblock-section-request-body, not in
  // the Parameters table. Grab the textarea that Swagger renders for the
  // "Example Value" / try-it form (best-effort, schema-unaware).
  if (!result.body) {
    const ta = opblock.querySelector<HTMLTextAreaElement>(
      ".opblock-section-request-body textarea, .body-param textarea"
    );
    const ex = ta?.value?.trim();
    if (ex && (ex.startsWith("{") || ex.startsWith("["))) {
      try { result.body = JSON.parse(ex); }
      catch { result.body = ex; }
    }
  }

  return result;
}

function mergePathParams(fromUrl: Param[], scraped: Param[]): Param[] {
  const byName = new Map<string, Param>();
  for (const p of fromUrl) byName.set(p.name, { ...p });
  for (const p of scraped) {
    const existing = byName.get(p.name);
    if (existing) Object.assign(existing, { ...existing, ...p });
    else byName.set(p.name, p);
  }
  return Array.from(byName.values());
}

// inferBody picks an encoding for the body slot:
//   - scraped JSON body (OAS3 application/json, or OAS2 body param) → json
//   - scraped formData rows → form
//   - otherwise → none (the user can change it in the editor)
function inferBody(method: string, json: unknown | null, formData: Param[]): BodySpec {
  if (json !== null) return { encoding: "json", json } as BodySpec;
  if (formData.length > 0) return { encoding: "form", form: formData };
  if (method === "GET" || method === "HEAD" || method === "DELETE") return { encoding: "none" };
  return { encoding: "none" };
}

function detectBaseUrl(): string | null {
  // 1. OpenAPI 3.x: Swagger UI renders a <select> with each `servers[].url`.
  const sel = document.querySelector<HTMLSelectElement>(
    ".servers select, select[aria-label='Servers']"
  );
  if (sel?.value) return trimSlash(sel.value);

  // 2. Swagger 2.0 (OpenAPI 2): no servers dropdown. The host+basePath is
  //    rendered into .info as text like "[ Base URL: petstore.swagger.io/v2 ]".
  //    The scheme comes from .scheme-container select (HTTP/HTTPS dropdown).
  const info = document.querySelector(".info");
  if (info) {
    const text = info.textContent ?? "";
    // Tolerate "Base URL: X", "Base URL : X", and optional surrounding [ ].
    const m = /Base\s*URL\s*[:：]\s*([^\s\]]+)/i.exec(text);
    if (m && m[1] && !/^\d+\.\d+/.test(m[1])) {
      let host = m[1].trim();
      if (!/^https?:\/\//i.test(host)) {
        host = pickScheme() + "://" + host.replace(/^\/+/, "");
      }
      return trimSlash(host);
    }
  }

  // 3. Last resort: page origin. The user can override in the panel.
  return trimSlash(window.location.origin);
}

function pickScheme(): string {
  const schemeSel = document.querySelector<HTMLSelectElement>(
    ".scheme-container select, select[aria-label='Schemes']"
  );
  const v = schemeSel?.value?.trim().toLowerCase();
  if (v === "http" || v === "https") return v;
  // Fall back to the docs page's own scheme (almost always https).
  return window.location.protocol.replace(":", "") || "https";
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

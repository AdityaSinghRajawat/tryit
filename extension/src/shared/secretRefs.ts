// Client-side mirror of internal/model/requestSpec.go::SecretRefs.
// Scans auth.valueRef / auth.username / auth.password for {{secret:NAME}}
// placeholders. Order-preserving; deduplicated.

import type { RequestSpec } from "./types";

const PLACEHOLDER_RE = /\{\{\s*secret\s*:\s*([A-Z0-9_]+)\s*\}\}/g;

export function extractSecretRefs(spec: RequestSpec): string[] {
  const seen = new Set<string>();
  const out: string[] = [];
  const sources = [spec.auth.valueRef, spec.auth.username, spec.auth.password];
  for (const src of sources) {
    if (!src) continue;
    PLACEHOLDER_RE.lastIndex = 0;
    let m: RegExpExecArray | null;
    while ((m = PLACEHOLDER_RE.exec(src)) !== null) {
      const name = m[1];
      if (seen.has(name)) continue;
      seen.add(name);
      out.push(name);
    }
  }
  return out;
}

export function hostOf(baseUrl: string): string {
  try {
    return new URL(baseUrl).host;
  } catch {
    return "";
  }
}

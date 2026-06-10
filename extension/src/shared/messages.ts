// IMPL §6.2 — content script → background → panel.
//
// Phase 1 carries a structuredHint produced by the Swagger detector and
// converted into a RequestSpec inside the panel (the no-AI vertical slice;
// Phase 2 moves the conversion server-side behind /parse).

import type { RequestSpec } from "./types";

export type ContentMsg = {
  kind: "tryit";
  pageUrl: string;
  scopedMarkdown: string;
  authSectionMarkdown?: string;
  framework?: "swagger" | "redoc" | "generic";
  structuredHint?: unknown;
};

// Phase 1 only: the content script already has enough information from the
// Swagger DOM to compute a RequestSpec; it includes it as a hint so the panel
// can render the editable form immediately, with no server round-trip.
export type Phase1Hint = {
  requestSpec: RequestSpec;
};

export type PanelMsg = ContentMsg;

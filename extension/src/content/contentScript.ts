// Content script orchestrator (isolated world).
//
// Swagger UI → client-side hint (Phase 1 fast path; no AI round-trip).
// Redoc / generic → scopedMarkdown + auth-section markdown (Phase 2 cascade).
// The background opens the side panel and forwards the message verbatim.

import { detect, frameworks, type Framework } from "./endpointDetector";
import { injectButton } from "./buttonInjector";
import type { ContentMsg } from "../shared/messages";

let activeFramework: Framework = "none";

function main(): void {
  activeFramework = detect();
  if (activeFramework === "none") return;

  scanAndInject();

  // Doc sites lazy-render (Swagger group toggle, Redoc lazy ops, generic SPA
  // routing) — observe and re-scan.
  const obs = new MutationObserver(() => scanAndInject());
  obs.observe(document.body, { subtree: true, childList: true });
}

function scanAndInject(): void {
  switch (activeFramework) {
    case "swagger":
      for (const block of frameworks.swagger.findEndpointBlocks()) {
        const anchor = frameworks.swagger.findInjectionAnchor(block);
        if (anchor) injectButton(anchor, () => onSwaggerTryIt(block));
      }
      return;
    case "redoc":
      for (const block of frameworks.redoc.findEndpointBlocks()) {
        const anchor = frameworks.redoc.findInjectionAnchor(block);
        if (anchor) injectButton(anchor, () => onProseTryIt(block, "redoc"));
      }
      return;
    case "generic":
      for (const block of frameworks.generic.findEndpointBlocks()) {
        const anchor = frameworks.generic.findInjectionAnchor(block);
        if (anchor) injectButton(anchor, () => onProseTryIt(block, "generic"));
      }
      return;
  }
}

function onSwaggerTryIt(block: HTMLElement): void {
  const spec = frameworks.swagger.buildRequestSpec(block);
  if (!spec) {
    console.warn("tryit: could not extract a RequestSpec from this opblock");
    return;
  }
  send({
    kind: "tryit",
    pageUrl: window.location.href,
    scopedMarkdown: "",
    framework: "swagger",
    structuredHint: { requestSpec: spec },
  });
}

function onProseTryIt(block: HTMLElement, framework: "redoc" | "generic"): void {
  const scope = frameworks[framework].scopeMarkdown(block);
  send({
    kind: "tryit",
    pageUrl: window.location.href,
    scopedMarkdown: scope.scopedMarkdown,
    authSectionMarkdown: scope.authSectionMarkdown || undefined,
    framework,
  });
}

function send(msg: ContentMsg): void {
  chrome.runtime.sendMessage(msg).catch((err) => {
    console.warn("tryit: sendMessage failed", err);
  });
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", main, { once: true });
} else {
  main();
}

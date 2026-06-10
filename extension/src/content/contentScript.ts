// Content script orchestrator (isolated world).
//
// Phase 1: detect Swagger UI → for each .opblock, inject a closed-shadow-DOM
// "Try it" button → on click, build a RequestSpec hint client-side → post to
// background, which opens the side panel and stashes the message for the
// panel to read.

import { detect, frameworks } from "./endpointDetector";
import { injectButton } from "./buttonInjector";
import type { ContentMsg } from "../shared/messages";

function main(): void {
  const framework = detect();
  if (framework === "none") return;

  // Initial pass.
  scanAndInject();

  // Swagger UI lazy-renders opblocks on group toggle; observe.
  const obs = new MutationObserver(() => scanAndInject());
  obs.observe(document.body, { subtree: true, childList: true });
}

function scanAndInject(): void {
  for (const block of frameworks.swagger.findEndpointBlocks()) {
    const anchor = frameworks.swagger.findInjectionAnchor(block);
    if (!anchor) continue;
    injectButton(anchor, () => onTryIt(block));
  }
}

function onTryIt(block: HTMLElement): void {
  const spec = frameworks.swagger.buildRequestSpec(block);
  if (!spec) {
    console.warn("tryit: could not extract a RequestSpec from this opblock");
    return;
  }
  const msg: ContentMsg = {
    kind: "tryit",
    pageUrl: window.location.href,
    scopedMarkdown: "",
    framework: "swagger",
    structuredHint: { requestSpec: spec },
  };
  chrome.runtime.sendMessage(msg).catch((err) => {
    console.warn("tryit: sendMessage failed", err);
  });
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", main, { once: true });
} else {
  main();
}

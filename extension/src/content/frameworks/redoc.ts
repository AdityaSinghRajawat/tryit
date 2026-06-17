// Redoc detection + scoping. Redoc renders each operation as a section
// with id="operation/<operationId>"; the heading is the H2/H3 inside.
// We don't try to harvest an OpenAPI hint client-side — the AI cascade
// gets the scoped markdown for the operation block and the page-level
// auth section.

import { elementToMarkdown } from "../markdownConverter";
import { findAuthSection } from "../domScoper";

const OPERATION_SELECTOR = "[data-section-id^='operation/'], section[id^='operation/']";

export function isRedoc(): boolean {
  if (document.querySelector("redoc, .redoc-wrap")) return true;
  if (document.querySelector(OPERATION_SELECTOR)) return true;
  return false;
}

export function findEndpointBlocks(): HTMLElement[] {
  return Array.from(document.querySelectorAll<HTMLElement>(OPERATION_SELECTOR));
}

export function findInjectionAnchor(block: HTMLElement): HTMLElement | null {
  return (
    block.querySelector<HTMLElement>("h1, h2, h3, h4") ??
    block.querySelector<HTMLElement>("[role='heading']") ??
    block
  );
}

export function scopeMarkdown(block: HTMLElement): {
  scopedMarkdown: string;
  authSectionMarkdown: string;
} {
  return {
    scopedMarkdown: elementToMarkdown(block),
    authSectionMarkdown: elementToMarkdown(findAuthSection()),
  };
}

// Generic heuristic detector — IMPL §9.3.
//
// A region candidate is a heading where the heading text OR the next code/
// pre block contains an HTTP method (GET|POST|...|OPTIONS) AND a URL or
// /path. Each match becomes an endpoint block; the button injects on the
// heading itself; markdown scoping uses the heading→next-heading walk in
// domScoper.

import { elementToMarkdown } from "../markdownConverter";
import { findAuthSection, scopeEndpointRegion } from "../domScoper";

const HEADING_TAGS = ["H1", "H2", "H3", "H4", "H5", "H6"];
const METHOD_RE = /\b(GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS)\b/;
const URL_OR_PATH_RE = /(https?:\/\/\S+|\/[A-Za-z0-9._~\-/{}:]+)/;

export function findEndpointBlocks(): HTMLElement[] {
  const out: HTMLElement[] = [];
  for (const tag of HEADING_TAGS) {
    for (const h of Array.from(document.getElementsByTagName(tag)) as HTMLElement[]) {
      if (looksLikeEndpoint(h)) out.push(h);
    }
  }
  return out;
}

export function findInjectionAnchor(block: HTMLElement): HTMLElement | null {
  return block;
}

export function scopeMarkdown(block: HTMLElement): {
  scopedMarkdown: string;
  authSectionMarkdown: string;
} {
  return {
    scopedMarkdown: elementToMarkdown(scopeEndpointRegion(block)),
    authSectionMarkdown: elementToMarkdown(findAuthSection()),
  };
}

function looksLikeEndpoint(heading: HTMLElement): boolean {
  const headingText = (heading.textContent ?? "").trim();
  if (hasMethodAndTarget(headingText)) return true;

  const next = heading.nextElementSibling;
  if (!next) return false;
  if (!isCodeOrPre(next)) return false;
  return hasMethodAndTarget((next.textContent ?? "").trim());
}

function hasMethodAndTarget(text: string): boolean {
  if (!text) return false;
  return METHOD_RE.test(text) && URL_OR_PATH_RE.test(text);
}

function isCodeOrPre(el: Element): boolean {
  return el.tagName === "PRE" || el.tagName === "CODE" || !!el.querySelector("pre, code");
}

// Turndown-backed HTML→markdown for the Phase 2 prose cascade. Strips
// non-content chrome (script/style/nav/footer/aside/form/iframe/button)
// before conversion so the AI prompt stays small and on-topic.

import TurndownService from "turndown";

const STRIP_SELECTORS = [
  "script",
  "style",
  "noscript",
  "iframe",
  "nav",
  "footer",
  "aside",
  "form",
  "button",
  "input",
  "select",
  "textarea",
  "[aria-hidden='true']",
];

let svc: TurndownService | null = null;

function turndown(): TurndownService {
  if (svc) return svc;
  svc = new TurndownService({
    headingStyle: "atx",
    codeBlockStyle: "fenced",
    fence: "```",
    bulletListMarker: "-",
    emDelimiter: "_",
  });
  svc.addRule("ignoreNoise", { filter: ["script", "style", "noscript"], replacement: () => "" });
  return svc;
}

export function htmlToMarkdown(html: string): string {
  if (!html) return "";
  return turndown().turndown(html);
}

export function elementToMarkdown(el: Element | null | undefined): string {
  if (!el) return "";
  return htmlToMarkdown(cleanHtml(el));
}

function cleanHtml(el: Element): string {
  const clone = el.cloneNode(true) as Element;
  for (const sel of STRIP_SELECTORS) {
    clone.querySelectorAll(sel).forEach((n) => n.remove());
  }
  return clone.outerHTML;
}

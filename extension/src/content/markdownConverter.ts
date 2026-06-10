// Phase 1 stub: structured detection (Swagger) doesn't need markdown
// extraction. Phase 2 wires up Turndown for prose docs.

import TurndownService from "turndown";

let svc: TurndownService | null = null;

export function htmlToMarkdown(html: string): string {
  if (!html) return "";
  svc ??= new TurndownService({ headingStyle: "atx", codeBlockStyle: "fenced" });
  return svc.turndown(html);
}

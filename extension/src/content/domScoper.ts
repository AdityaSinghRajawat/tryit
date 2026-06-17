// Endpoint + authentication region scoping per IMPL §9.3.
//
// scopeEndpointRegion walks from a heading element to the next sibling
// heading of the same or higher level, returning the contiguous element
// range as a synthetic container. findAuthSection scans the page for a
// heading whose text matches the auth vocabulary and scopes the same way.

const HEADING_TAGS = ["H1", "H2", "H3", "H4", "H5", "H6"];
const AUTH_HEADING_RE = /^\s*(authentication|authorization|auth|credentials?|tokens?|api\s*keys?)\b/i;

export function isHeading(el: Element | null | undefined): el is HTMLElement {
  return !!el && HEADING_TAGS.includes(el.tagName);
}

export function headingLevel(el: Element): number {
  const m = /^H([1-6])$/.exec(el.tagName);
  return m ? Number(m[1]) : 99;
}

// scopeEndpointRegion collects [heading, ...siblings until next heading of
// equal/greater level]. Returns null if `start` isn't a heading-rooted
// region. The caller passes the result to elementToMarkdown.
export function scopeEndpointRegion(start: Element): HTMLElement | null {
  if (!isHeading(start)) return null;
  const level = headingLevel(start);
  const container = document.createElement("div");
  container.appendChild(start.cloneNode(true));

  let cursor: Element | null = start.nextElementSibling;
  while (cursor) {
    if (isHeading(cursor) && headingLevel(cursor) <= level) break;
    container.appendChild(cursor.cloneNode(true));
    cursor = cursor.nextElementSibling;
  }
  return container;
}

// findAuthSection returns the page-level Authentication section, or null
// if no matching heading exists. Match is case-insensitive on the heading
// text only.
export function findAuthSection(): HTMLElement | null {
  for (const tag of HEADING_TAGS) {
    const headings = document.getElementsByTagName(tag);
    for (const h of Array.from(headings)) {
      const text = (h.textContent ?? "").trim();
      if (AUTH_HEADING_RE.test(text)) {
        return scopeEndpointRegion(h);
      }
    }
  }
  return null;
}

// D13: background is a RELAY ONLY. It opens the panel and stashes the
// triggering message so the panel can pick it up. It never calls the Go
// server. The panel does all server calls directly via serverClient.ts.

import type { ContentMsg } from "../shared/messages";
import { RUNTIME_TRIGGER, STORAGE_KEY_LAST_TRIGGER } from "../shared/constants";

// Allow the action button to open the panel even with no detected endpoint.
chrome.sidePanel
  .setPanelBehavior({ openPanelOnActionClick: true })
  .catch((err) => console.warn("setPanelBehavior failed:", err));

chrome.runtime.onMessage.addListener((msg: unknown, sender, sendResponse) => {
  if (!isContentMsg(msg)) return false;

  // Stash for the panel; the panel reads on mount.
  chrome.storage.local
    .set({ [STORAGE_KEY_LAST_TRIGGER]: msg })
    .catch((err) => console.warn("storage.local.set failed:", err));

  const tabId = sender.tab?.id;
  if (tabId !== undefined) {
    chrome.sidePanel
      .open({ tabId })
      .catch((err) => console.warn("sidePanel.open failed:", err));
  }

  // Also broadcast a runtime message: if the panel is already open it can
  // skip the storage read.
  chrome.runtime.sendMessage({ kind: RUNTIME_TRIGGER, payload: msg }).catch(() => {
    /* no listeners is fine */
  });

  sendResponse({ ok: true });
  return true;
});

function isContentMsg(x: unknown): x is ContentMsg {
  return !!x && typeof x === "object" && (x as { kind?: string }).kind === "tryit";
}

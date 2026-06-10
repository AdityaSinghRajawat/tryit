import { useCallback, useEffect, useState } from "react";
import { PairDialog } from "./components/PairDialog";
import { RequestEditor } from "./components/RequestEditor";
import { ResponseView } from "./components/ResponseView";
import { ServerError, serverClient } from "./serverClient";
import {
  STORAGE_KEY_LAST_TRIGGER,
} from "../shared/constants";
import type { ContentMsg, Phase1Hint } from "../shared/messages";
import type { ExecuteResponse, RequestSpec } from "../shared/types";

type Stage =
  | { kind: "boot" }
  | { kind: "pair-needed" }
  | { kind: "idle" }
  | { kind: "request"; spec: RequestSpec }
  | { kind: "sent"; resp: ExecuteResponse; spec: RequestSpec };

export function App() {
  const [stage, setStage] = useState<Stage>({ kind: "boot" });
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const refreshFromTrigger = useCallback(async () => {
    const r = await chrome.storage.local.get(STORAGE_KEY_LAST_TRIGGER);
    const msg = r[STORAGE_KEY_LAST_TRIGGER] as ContentMsg | undefined;
    if (msg && msg.kind === "tryit") {
      const hint = msg.structuredHint as Phase1Hint | undefined;
      if (hint?.requestSpec) {
        setStage({ kind: "request", spec: hint.requestSpec });
        // Clear the trigger so reload doesn't reuse stale state.
        await chrome.storage.local.remove(STORAGE_KEY_LAST_TRIGGER);
      }
    }
  }, []);

  // Boot: check pairing, then load any pending trigger.
  useEffect(() => {
    (async () => {
      const has = await serverClient.hasToken();
      if (!has) {
        setStage({ kind: "pair-needed" });
        return;
      }
      try {
        await serverClient.health();
      } catch (e) {
        // Server unreachable: still allow the panel to render the editor;
        // sending will surface the error.
        console.warn("health check failed:", e);
      }
      setStage({ kind: "idle" });
      await refreshFromTrigger();
    })();
  }, [refreshFromTrigger]);

  // Listen for new triggers while open.
  useEffect(() => {
    const onChange = (changes: Record<string, chrome.storage.StorageChange>) => {
      if (STORAGE_KEY_LAST_TRIGGER in changes && changes[STORAGE_KEY_LAST_TRIGGER].newValue) {
        refreshFromTrigger();
      }
    };
    chrome.storage.local.onChanged.addListener(onChange);
    return () => chrome.storage.local.onChanged.removeListener(onChange);
  }, [refreshFromTrigger]);

  const send = async (spec: RequestSpec) => {
    setBusy(true);
    setError(null);
    try {
      const resp = await serverClient.execute({ requestSpec: spec });
      setStage({ kind: "sent", resp, spec });
    } catch (e) {
      const msg = e instanceof ServerError ? `${e.code}: ${e.message}` : String(e);
      setError(msg);
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="app">
      <Style />
      <header className="topbar">
        <strong>tryit</strong>
        <span className="phase">Phase 1</span>
      </header>

      {stage.kind === "boot" && <p className="loading">Loading…</p>}

      {stage.kind === "pair-needed" && (
        <PairDialog onPaired={() => setStage({ kind: "idle" })} />
      )}

      {stage.kind === "idle" && (
        <div className="idle">
          <p>
            Visit a Swagger UI page (e.g. <code>petstore.swagger.io</code>) and
            click the <strong>Try it</strong> button next to an endpoint to
            load it here.
          </p>
        </div>
      )}

      {stage.kind === "request" && (
        <RequestEditor initial={stage.spec} busy={busy} onSend={send} />
      )}

      {stage.kind === "sent" && (
        <>
          <RequestEditor initial={stage.spec} busy={busy} onSend={send} />
          <ResponseView resp={stage.resp} />
        </>
      )}

      {error && <p className="error">{error}</p>}
    </div>
  );
}

function Style() {
  return (
    <style>{`
      body, html, #root { margin:0; padding:0; height:100%; font: 13px/1.4 system-ui, sans-serif; color:#222; background:#fafafa; }
      .app { padding: 12px; display:flex; flex-direction:column; gap:12px; }
      .topbar { display:flex; align-items:baseline; gap:8px; }
      .topbar .phase { font-size: 11px; color:#888; }
      .dialog h2 { margin:0 0 8px; font-size: 15px; }
      .dialog .hint { color:#555; margin: 0 0 12px; }
      label { display:block; font-size: 11px; color:#555; }
      input, select, textarea {
        width: 100%; box-sizing: border-box; font: inherit; padding: 6px 8px;
        border: 1px solid #d0d0d0; border-radius: 6px; background:#fff;
      }
      .row { margin-bottom: 8px; }
      .block { border-top: 1px solid #eee; padding-top: 10px; margin-top: 10px; }
      .block h3 { font-size: 12px; text-transform: uppercase; color:#666; margin: 0 0 8px; letter-spacing: 0.04em; }
      button {
        font: 600 13px/1 system-ui; padding: 8px 12px; border-radius: 6px; border: 1px solid #2a78c0;
        background: #2a78c0; color:white; cursor:pointer;
      }
      button:disabled { opacity: 0.5; cursor: default; }
      .send-row { margin-top: 8px; }
      .response { border-top: 1px solid #eee; padding-top: 10px; }
      .badge { padding: 2px 6px; border-radius: 4px; font-weight: 600; }
      .status-2xx { background:#d6f5d6; color:#1e6b1e; }
      .status-3xx { background:#fff1c0; color:#7a5d00; }
      .status-4xx, .status-5xx { background:#fbd6d6; color:#8b1212; }
      .duration { margin-left: 8px; color:#666; }
      .truncated { margin-left: 8px; color:#8b1212; font-weight: 600; }
      pre { background:#f3f3f3; padding: 8px; border-radius: 6px; overflow:auto; white-space: pre-wrap; word-break: break-all; }
      .error { color:#8b1212; }
      code { background:#eee; padding: 1px 4px; border-radius: 3px; }
    `}</style>
  );
}

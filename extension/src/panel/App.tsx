import { useCallback, useEffect, useState } from "react";
import { PairDialog } from "./components/PairDialog";
import { RequestEditor } from "./components/RequestEditor";
import { ResponseView } from "./components/ResponseView";
import { ConfirmBar } from "./components/ConfirmBar";
import { SecretMapper } from "./components/SecretMapper";
import { ConsentDialog } from "./components/ConsentDialog";
import { CodeGenView } from "./components/CodeGenView";
import { ServerError, serverClient } from "./serverClient";
import { STORAGE_KEY_LAST_TRIGGER } from "../shared/constants";
import { extractSecretRefs, hostOf } from "../shared/secretRefs";
import type { ContentMsg, Phase1Hint } from "../shared/messages";
import type {
  ExecuteResponse,
  ParseRequest,
  ParseResponse,
  RequestSpec,
} from "../shared/types";

type Stage =
  | { kind: "boot" }
  | { kind: "pair-needed" }
  | { kind: "idle" }
  | { kind: "parsing"; parseInput: ParseRequest }
  | { kind: "request"; spec: RequestSpec; parse?: ParseResponse; parseInput?: ParseRequest }
  | { kind: "mapping"; spec: RequestSpec; parse?: ParseResponse; parseInput?: ParseRequest }
  | {
      kind: "consent";
      spec: RequestSpec;
      refs: Record<string, string>;
      required: { secret: string; host: string };
      parseInput?: ParseRequest;
    }
  | { kind: "sent"; spec: RequestSpec; resp: ExecuteResponse; parseInput?: ParseRequest };

export function App() {
  const [stage, setStage] = useState<Stage>({ kind: "boot" });
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const consumeTrigger = useCallback(async () => {
    const r = await chrome.storage.local.get(STORAGE_KEY_LAST_TRIGGER);
    const msg = r[STORAGE_KEY_LAST_TRIGGER] as ContentMsg | undefined;
    if (!msg || msg.kind !== "tryit") return;
    await chrome.storage.local.remove(STORAGE_KEY_LAST_TRIGGER);

    const hint = msg.structuredHint as Phase1Hint | undefined;
    if (hint?.requestSpec) {
      setStage({ kind: "request", spec: hint.requestSpec });
      return;
    }
    if (msg.scopedMarkdown) {
      const parseInput: ParseRequest = {
        pageUrl: msg.pageUrl,
        scopedMarkdown: msg.scopedMarkdown,
        authSectionMarkdown: msg.authSectionMarkdown,
        framework: msg.framework,
      };
      setStage({ kind: "parsing", parseInput });
    }
  }, []);

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
        console.warn("health check failed:", e);
      }
      setStage({ kind: "idle" });
      await consumeTrigger();
    })();
  }, [consumeTrigger]);

  useEffect(() => {
    const onChange = (changes: Record<string, chrome.storage.StorageChange>) => {
      if (STORAGE_KEY_LAST_TRIGGER in changes && changes[STORAGE_KEY_LAST_TRIGGER].newValue) {
        consumeTrigger();
      }
    };
    chrome.storage.local.onChanged.addListener(onChange);
    return () => chrome.storage.local.onChanged.removeListener(onChange);
  }, [consumeTrigger]);

  // Parse cascade: parsing → request.
  useEffect(() => {
    if (stage.kind !== "parsing") return;
    let cancelled = false;
    (async () => {
      setError(null);
      try {
        const parse = await serverClient.parse(stage.parseInput);
        if (cancelled) return;
        setStage({
          kind: "request",
          spec: parse.requestSpec,
          parse,
          parseInput: stage.parseInput,
        });
      } catch (err) {
        if (cancelled) return;
        setError(formatError(err));
        setStage({ kind: "idle" });
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [stage]);

  const onSendRequest = (spec: RequestSpec) => {
    if (stage.kind !== "request") return;
    if (extractSecretRefs(spec).length === 0) {
      void doExecute(spec, {});
      return;
    }
    setStage({ ...stage, kind: "mapping", spec });
  };

  const onMappingReady = (refs: Record<string, string>) => {
    if (stage.kind !== "mapping") return;
    void doExecute(stage.spec, refs);
  };

  const onConsentGranted = () => {
    if (stage.kind !== "consent") return;
    void doExecute(stage.spec, stage.refs);
  };

  const doExecute = async (spec: RequestSpec, refs: Record<string, string>) => {
    setBusy(true);
    setError(null);
    try {
      const resp = await serverClient.execute({ requestSpec: spec, secretRefs: refs });
      if (resp.consentRequired) {
        setStage((prev) => ({
          kind: "consent",
          spec,
          refs,
          required: resp.consentRequired!,
          parseInput: anyParseInput(prev),
        }));
        return;
      }
      setStage((prev) => ({
        kind: "sent",
        spec,
        resp,
        parseInput: anyParseInput(prev),
      }));
    } catch (err) {
      setError(formatError(err));
    } finally {
      setBusy(false);
    }
  };

  const reparse = () => {
    const parseInput = anyParseInput(stage);
    if (!parseInput) return;
    setStage({ kind: "parsing", parseInput: { ...parseInput, force: true } });
  };

  const dismissConfirmBar = () => {
    if (stage.kind !== "request") return;
    setStage({ ...stage, parse: undefined });
  };

  return (
    <div className="app">
      <Style />
      <header className="topbar">
        <strong>tryit</strong>
        <span className="phase">Phase 3</span>
      </header>

      {stage.kind === "boot" && <p className="loading">Loading…</p>}

      {stage.kind === "pair-needed" && (
        <PairDialog onPaired={() => setStage({ kind: "idle" })} />
      )}

      {stage.kind === "idle" && (
        <div className="idle">
          <p>
            Open an API docs page (Swagger, Redoc, or a prose doc) and click the
            <strong> Try it</strong> button next to an endpoint.
          </p>
        </div>
      )}

      {stage.kind === "parsing" && (
        <p className="loading">Parsing endpoint with the AI cascade…</p>
      )}

      {stage.kind === "request" && (
        <>
          {stage.parse?.needsConfirmation && (
            <ConfirmBar
              source={stage.parse.source}
              confidence={stage.parse.confidence}
              notes={stage.spec.notes}
              onConfirm={dismissConfirmBar}
              onReparse={reparse}
            />
          )}
          <RequestEditor initial={stage.spec} busy={busy} onSend={onSendRequest} />
          <HostBadge baseUrl={stage.spec.baseUrl} />
        </>
      )}

      {stage.kind === "mapping" && (
        <SecretMapper
          spec={stage.spec}
          onReady={onMappingReady}
          onCancel={() =>
            setStage({
              kind: "request",
              spec: stage.spec,
              parse: stage.parse,
              parseInput: stage.parseInput,
            })
          }
        />
      )}

      {stage.kind === "consent" && (
        <ConsentDialog
          secret={stage.required.secret}
          host={stage.required.host}
          onGranted={onConsentGranted}
          onCancel={() =>
            setStage({
              kind: "request",
              spec: stage.spec,
              parseInput: stage.parseInput,
            })
          }
        />
      )}

      {stage.kind === "sent" && (
        <>
          <RequestEditor initial={stage.spec} busy={busy} onSend={onSendRequest} />
          <ResponseView resp={stage.resp} />
          <CodeGenView spec={stage.spec} />
        </>
      )}

      {error && <p className="error">{error}</p>}
    </div>
  );
}

function HostBadge({ baseUrl }: { baseUrl: string }) {
  const host = hostOf(baseUrl);
  if (!host) return null;
  return (
    <p className="muted">
      Target host: <code>{host}</code>
    </p>
  );
}

function anyParseInput(s: Stage): ParseRequest | undefined {
  if ("parseInput" in s) return s.parseInput;
  return undefined;
}

function formatError(err: unknown): string {
  return err instanceof ServerError ? `${err.code}: ${err.message}` : String(err);
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
      .dialog-actions { display:flex; gap:8px; justify-content:flex-end; }
      label { display:block; font-size: 11px; color:#555; }
      input, select, textarea {
        width: 100%; box-sizing: border-box; font: inherit; padding: 6px 8px;
        border: 1px solid #d0d0d0; border-radius: 6px; background:#fff;
      }
      .row { margin-bottom: 8px; }
      .block { border-top: 1px solid #eee; padding-top: 10px; margin-top: 10px; }
      .block h3 { font-size: 12px; text-transform: uppercase; color:#666; margin: 0 0 8px; letter-spacing: 0.04em; display:flex; align-items:center; justify-content:space-between; }
      .muted { color:#888; font-size: 11px; margin: 0 0 6px; }
      .warn { display:block; color:#b56300; font-size: 11px; margin-top: 4px; }
      textarea { font-family: ui-monospace, SFMono-Regular, monospace; font-size: 12px; }
      .kv-row { display:flex; gap:6px; margin-bottom: 6px; align-items:center; }
      .kv-row .kv-name { flex: 0 0 35%; }
      .kv-row .kv-value { flex: 1 1 auto; }
      .kv-row .del {
        flex: 0 0 auto; padding: 2px 8px; background:transparent; color:#888;
        border:1px solid #d0d0d0; border-radius:6px;
      }
      .kv-row .del:hover { color:#8b1212; border-color:#c88; }
      .add {
        font: 500 11px/1 system-ui; padding: 3px 8px; border-radius: 4px;
        border: 1px solid #d0d0d0; background:#fff; color:#444; cursor:pointer;
      }
      .add:hover { background:#f4f4f4; }
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
      .loading { color:#666; font-style: italic; }

      .confirm-bar {
        background: #fff8df; border: 1px solid #e9d77a; border-radius: 6px;
        padding: 10px 12px;
      }
      .confirm-header { display:flex; justify-content:space-between; align-items:center; gap:8px; }
      .confirm-pill { font-size: 11px; font-weight: 600; padding: 2px 8px; border-radius: 999px; }
      .confirm-low  { background:#fbd6d6; color:#8b1212; }
      .confirm-mid  { background:#fff1c0; color:#7a5d00; }
      .confirm-high { background:#d6f5d6; color:#1e6b1e; }
      .confirm-actions { display:flex; gap:6px; }
      .confirm-notes { margin-top: 8px; }
      .confirm-help  { margin-top: 4px; }

      .codegen .codegen-tabs { display:inline-flex; gap:6px; align-items:center; }
      .codegen-tab.active { background:#2a78c0; color:white; border-color:#2a78c0; }
      .codegen-snippet { max-height: 320px; }
    `}</style>
  );
}

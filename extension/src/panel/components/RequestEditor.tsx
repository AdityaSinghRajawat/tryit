import { useEffect, useState } from "react";
import type { RequestSpec } from "../../shared/types";

interface Props {
  initial: RequestSpec;
  busy: boolean;
  onSend: (spec: RequestSpec) => void;
}

// Minimal Phase 1 editor: lets the user adjust method/baseUrl/path and edit
// path-param values inline. Phase 2 grows full body/query/header/auth UIs.
export function RequestEditor({ initial, busy, onSend }: Props) {
  const [spec, setSpec] = useState<RequestSpec>(initial);

  useEffect(() => {
    setSpec(initial);
  }, [initial]);

  const setPathParam = (idx: number, value: string) => {
    setSpec((s) => {
      const next = { ...s, pathParams: [...(s.pathParams ?? [])] };
      next.pathParams[idx] = { ...next.pathParams[idx], value };
      return next;
    });
  };

  const setAuth = (patch: Partial<RequestSpec["auth"]>) => {
    setSpec((s) => ({ ...s, auth: { ...s.auth, ...patch } }));
  };

  return (
    <div className="editor">
      <div className="row">
        <label>
          Method
          <select
            value={spec.method}
            onChange={(e) => setSpec((s) => ({ ...s, method: e.target.value as RequestSpec["method"] }))}
          >
            {["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"].map((m) => (
              <option key={m} value={m}>
                {m}
              </option>
            ))}
          </select>
        </label>
      </div>

      <div className="row">
        <label>
          Base URL
          <input
            value={spec.baseUrl}
            onChange={(e) => setSpec((s) => ({ ...s, baseUrl: e.target.value }))}
            placeholder="https://api.example.com"
            spellCheck={false}
          />
        </label>
      </div>

      <div className="row">
        <label>
          Path
          <input
            value={spec.path}
            onChange={(e) => setSpec((s) => ({ ...s, path: e.target.value }))}
            placeholder="/v1/things/{id}"
            spellCheck={false}
          />
        </label>
      </div>

      {(spec.pathParams ?? []).length > 0 && (
        <div className="block">
          <h3>Path parameters</h3>
          {(spec.pathParams ?? []).map((p, i) => (
            <div className="row" key={p.name}>
              <label>
                {p.name}
                <input
                  value={p.value ?? ""}
                  onChange={(e) => setPathParam(i, e.target.value)}
                  placeholder={p.name}
                />
              </label>
            </div>
          ))}
        </div>
      )}

      <div className="block">
        <h3>Auth (Phase 1 — bearer/none only via env secret)</h3>
        <div className="row">
          <label>
            Type
            <select
              value={spec.auth.type}
              onChange={(e) =>
                setAuth({
                  type: e.target.value as RequestSpec["auth"]["type"],
                  valueRef: e.target.value === "bearer" ? "{{secret:DEFAULT}}" : undefined,
                })
              }
            >
              <option value="none">none</option>
              <option value="bearer">bearer</option>
            </select>
          </label>
        </div>
        {spec.auth.type === "bearer" && (
          <div className="row">
            <label>
              Secret name (TRYIT_SECRET_&lt;NAME&gt; env var)
              <input
                value={extractName(spec.auth.valueRef)}
                onChange={(e) => setAuth({ valueRef: `{{secret:${e.target.value.toUpperCase()}}}` })}
                placeholder="PETSTORE_KEY"
                spellCheck={false}
              />
            </label>
          </div>
        )}
      </div>

      <div className="row send-row">
        <button disabled={busy} onClick={() => onSend(spec)}>
          {busy ? "Sending…" : "Send"}
        </button>
      </div>
    </div>
  );
}

function extractName(template?: string): string {
  if (!template) return "";
  const m = /\{\{secret:([A-Z0-9_]+)\}\}/.exec(template);
  return m ? m[1] : "";
}

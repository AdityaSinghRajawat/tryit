import { useEffect, useState } from "react";
import type { BodySpec, Param, RequestSpec } from "../../shared/types";

interface Props {
  initial: RequestSpec;
  busy: boolean;
  onSend: (spec: RequestSpec) => void;
}

export function RequestEditor({ initial, busy, onSend }: Props) {
  const [spec, setSpec] = useState<RequestSpec>(initial);

  useEffect(() => {
    setSpec(normalise(initial));
  }, [initial]);

  const patch = (p: Partial<RequestSpec>) => setSpec((s) => ({ ...s, ...p }));

  const setPathParam = (idx: number, value: string) => {
    setSpec((s) => {
      const next = [...(s.pathParams ?? [])];
      next[idx] = { ...next[idx], value };
      return { ...s, pathParams: next };
    });
  };

  const setQuery = (idx: number, field: keyof Param, value: string) => {
    setSpec((s) => {
      const next = [...(s.query ?? [])];
      next[idx] = { ...next[idx], [field]: value };
      return { ...s, query: next };
    });
  };

  const addQuery = () => patch({ query: [...(spec.query ?? []), { name: "", value: "" }] });
  const removeQuery = (i: number) =>
    setSpec((s) => ({ ...s, query: (s.query ?? []).filter((_, idx) => idx !== i) }));

  const setHeader = (idx: number, field: "name" | "value", value: string) => {
    setSpec((s) => {
      const next = [...(s.headers ?? [])];
      next[idx] = { ...next[idx], [field]: value };
      return { ...s, headers: next };
    });
  };
  const addHeader = () => patch({ headers: [...(spec.headers ?? []), { name: "", value: "" }] });
  const removeHeader = (i: number) =>
    setSpec((s) => ({ ...s, headers: (s.headers ?? []).filter((_, idx) => idx !== i) }));

  const setBody = (b: Partial<BodySpec>) => setSpec((s) => ({ ...s, body: { ...s.body, ...b } }));
  const setForm = (idx: number, field: keyof Param, value: string) => {
    setSpec((s) => {
      const next = [...(s.body.form ?? [])];
      next[idx] = { ...next[idx], [field]: value };
      return { ...s, body: { ...s.body, form: next } };
    });
  };
  const addForm = () =>
    setBody({ form: [...(spec.body.form ?? []), { name: "", value: "" }] });
  const removeForm = (i: number) =>
    setSpec((s) => ({
      ...s,
      body: { ...s.body, form: (s.body.form ?? []).filter((_, idx) => idx !== i) },
    }));

  const setAuth = (p: Partial<RequestSpec["auth"]>) => patch({ auth: { ...spec.auth, ...p } });

  const submit = () => onSend(toSendShape(spec));

  return (
    <div className="editor">
      <div className="row">
        <label>
          Method
          <select
            value={spec.method}
            onChange={(e) => patch({ method: e.target.value as RequestSpec["method"] })}
          >
            {["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"].map((m) => (
              <option key={m} value={m}>{m}</option>
            ))}
          </select>
        </label>
      </div>

      <div className="row">
        <label>
          Base URL
          <input
            value={spec.baseUrl}
            onChange={(e) => patch({ baseUrl: e.target.value })}
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
            onChange={(e) => patch({ path: e.target.value })}
            placeholder="/v1/things/{id}"
            spellCheck={false}
          />
        </label>
      </div>

      {(spec.pathParams ?? []).length > 0 && (
        <div className="block">
          <h3>Path parameters</h3>
          {(spec.pathParams ?? []).map((p, i) => (
            <div className="row" key={`pp-${i}-${p.name}`}>
              <label>
                {p.name}{p.required ? " *" : ""}
                <input
                  value={p.value ?? ""}
                  onChange={(e) => setPathParam(i, e.target.value)}
                  placeholder={p.description || p.name}
                />
              </label>
            </div>
          ))}
        </div>
      )}

      <div className="block">
        <h3>
          Query parameters
          <button type="button" className="add" onClick={addQuery}>+ add</button>
        </h3>
        {(spec.query ?? []).length === 0 && <p className="muted">none</p>}
        {(spec.query ?? []).map((p, i) => (
          <div className="kv-row" key={`q-${i}`}>
            <input
              className="kv-name"
              value={p.name}
              onChange={(e) => setQuery(i, "name", e.target.value)}
              placeholder="name"
            />
            <input
              className="kv-value"
              value={p.value ?? ""}
              onChange={(e) => setQuery(i, "value", e.target.value)}
              placeholder={p.description || "value"}
            />
            <button type="button" className="del" onClick={() => removeQuery(i)}>×</button>
          </div>
        ))}
      </div>

      <div className="block">
        <h3>
          Headers
          <button type="button" className="add" onClick={addHeader}>+ add</button>
        </h3>
        {(spec.headers ?? []).length === 0 && <p className="muted">none</p>}
        {(spec.headers ?? []).map((h, i) => (
          <div className="kv-row" key={`h-${i}`}>
            <input
              className="kv-name"
              value={h.name}
              onChange={(e) => setHeader(i, "name", e.target.value)}
              placeholder="Header-Name"
            />
            <input
              className="kv-value"
              value={h.value}
              onChange={(e) => setHeader(i, "value", e.target.value)}
              placeholder="value"
            />
            <button type="button" className="del" onClick={() => removeHeader(i)}>×</button>
          </div>
        ))}
      </div>

      <div className="block">
        <h3>Body</h3>
        <div className="row">
          <label>
            Encoding
            <select
              value={spec.body.encoding}
              onChange={(e) => setBody({ encoding: e.target.value as BodySpec["encoding"] })}
            >
              <option value="none">none</option>
              <option value="json">json</option>
              <option value="form">form</option>
              <option value="raw">raw</option>
            </select>
          </label>
        </div>
        {spec.body.encoding === "json" && (
          <div className="row">
            <label>
              JSON
              <textarea
                rows={6}
                value={jsonText(spec.body)}
                onChange={(e) => setBody({ json: e.target.value as unknown as never })}
                placeholder='{"key": "value"}'
                spellCheck={false}
              />
              {jsonError(spec.body) && (
                <span className="warn">⚠ invalid JSON — fix before sending</span>
              )}
            </label>
          </div>
        )}
        {spec.body.encoding === "raw" && (
          <div className="row">
            <label>
              Raw body
              <textarea
                rows={6}
                value={spec.body.raw ?? ""}
                onChange={(e) => setBody({ raw: e.target.value })}
                spellCheck={false}
              />
            </label>
            <label>
              Content-Type override
              <input
                value={spec.body.contentType ?? ""}
                onChange={(e) => setBody({ contentType: e.target.value })}
                placeholder="text/plain"
                spellCheck={false}
              />
            </label>
          </div>
        )}
        {spec.body.encoding === "form" && (
          <>
            <p className="muted">application/x-www-form-urlencoded</p>
            {(spec.body.form ?? []).length === 0 && <p className="muted">none</p>}
            {(spec.body.form ?? []).map((p, i) => (
              <div className="kv-row" key={`f-${i}`}>
                <input
                  className="kv-name"
                  value={p.name}
                  onChange={(e) => setForm(i, "name", e.target.value)}
                  placeholder="name"
                />
                <input
                  className="kv-value"
                  value={p.value ?? ""}
                  onChange={(e) => setForm(i, "value", e.target.value)}
                  placeholder="value"
                />
                <button type="button" className="del" onClick={() => removeForm(i)}>×</button>
              </div>
            ))}
            <button type="button" className="add" onClick={addForm}>+ add field</button>
          </>
        )}
      </div>

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
        <button disabled={busy || (spec.body.encoding === "json" && !!jsonError(spec.body))} onClick={submit}>
          {busy ? "Sending…" : "Send"}
        </button>
      </div>
    </div>
  );
}

// --- helpers -----------------------------------------------------------

// While the user types JSON we hold it as a string in body.json (a string is a
// valid JSON-schema "json" value); on Send we parse it back to a real object.
function jsonText(body: BodySpec): string {
  if (body.json === undefined || body.json === null) return "";
  if (typeof body.json === "string") return body.json;
  try {
    return JSON.stringify(body.json, null, 2);
  } catch {
    return String(body.json);
  }
}

function jsonError(body: BodySpec): string | null {
  if (body.encoding !== "json") return null;
  const txt = jsonText(body).trim();
  if (!txt) return null;
  try {
    JSON.parse(txt);
    return null;
  } catch (e) {
    return (e as Error).message;
  }
}

// normalise pre-fills empty arrays so the editor never crashes on undefined.
function normalise(s: RequestSpec): RequestSpec {
  return {
    ...s,
    pathParams: s.pathParams ?? [],
    query: s.query ?? [],
    headers: s.headers ?? [],
    body: s.body ?? { encoding: "none" },
  };
}

// toSendShape: turn the editor state into a server-ready RequestSpec
// (JSON body string → parsed JSON; drop empty rows).
function toSendShape(s: RequestSpec): RequestSpec {
  const clean = (arr: Param[] | undefined) =>
    (arr ?? []).filter((p) => p.name.trim() !== "");
  const out: RequestSpec = {
    ...s,
    pathParams: clean(s.pathParams),
    query: clean(s.query),
    headers: (s.headers ?? []).filter((h) => h.name.trim() !== ""),
  };
  if (s.body.encoding === "json") {
    const txt = jsonText(s.body).trim();
    out.body = { ...s.body, json: txt ? JSON.parse(txt) : null };
  } else if (s.body.encoding === "form") {
    out.body = { ...s.body, form: clean(s.body.form) };
  }
  return out;
}

function extractName(template?: string): string {
  if (!template) return "";
  const m = /\{\{secret:([A-Z0-9_]+)\}\}/.exec(template);
  return m ? m[1] : "";
}

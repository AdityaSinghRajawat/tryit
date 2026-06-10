import type { ExecuteResponse } from "../../shared/types";

interface Props {
  resp: ExecuteResponse;
}

export function ResponseView({ resp }: Props) {
  const isJson = (resp.responseHeaders["Content-Type"] ?? "").toLowerCase().includes("json");
  let body = resp.body;
  if (isJson) {
    try {
      body = JSON.stringify(JSON.parse(resp.body), null, 2);
    } catch {
      /* ignore */
    }
  }
  return (
    <div className="response">
      <div className="status">
        <span className={`badge status-${Math.floor(resp.status / 100)}xx`}>
          {resp.status}
        </span>
        <span className="duration">{resp.durationMs} ms</span>
        {resp.truncated && <span className="truncated">truncated</span>}
      </div>

      <h3>Request preview (masked)</h3>
      <pre className="preview">
        {resp.requestPreview.method} {resp.requestPreview.url}
        {"\n"}
        {Object.entries(resp.requestPreview.headers)
          .map(([k, v]) => `${k}: ${v}`)
          .join("\n")}
        {resp.requestPreview.body ? `\n\n${resp.requestPreview.body}` : ""}
      </pre>

      <h3>Response headers</h3>
      <pre className="headers">
        {Object.entries(resp.responseHeaders)
          .map(([k, v]) => `${k}: ${v}`)
          .join("\n")}
      </pre>

      <h3>Body</h3>
      <pre className="body">{body}</pre>
    </div>
  );
}

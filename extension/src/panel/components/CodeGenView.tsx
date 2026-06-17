// Code-generation panel shown beneath the response. Three language tabs;
// the server renders each one via /generate.

import { useEffect, useState } from "react";
import { ServerError, serverClient } from "../serverClient";
import type { GenerateLanguage, RequestSpec } from "../../shared/types";

interface Props {
  spec: RequestSpec;
}

const LANGUAGES: { id: GenerateLanguage; label: string }[] = [
  { id: "curl", label: "curl" },
  { id: "python", label: "Python" },
  { id: "javascript", label: "JavaScript" },
];

export function CodeGenView({ spec }: Props) {
  const [language, setLanguage] = useState<GenerateLanguage>("curl");
  const [code, setCode] = useState<string>("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    let cancelled = false;
    setBusy(true);
    setError(null);
    setCopied(false);
    serverClient
      .generate({ requestSpec: spec, language })
      .then((r) => {
        if (!cancelled) setCode(r.code);
      })
      .catch((err) => {
        if (cancelled) return;
        setError(err instanceof ServerError ? `${err.code}: ${err.message}` : String(err));
        setCode("");
      })
      .finally(() => {
        if (!cancelled) setBusy(false);
      });
    return () => {
      cancelled = true;
    };
  }, [spec, language]);

  async function copy() {
    if (!code) return;
    try {
      await navigator.clipboard.writeText(code);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch (err) {
      setError("Copy failed: " + String(err));
    }
  }

  return (
    <div className="block codegen">
      <h3>
        Code
        <span className="codegen-tabs">
          {LANGUAGES.map((l) => (
            <button
              key={l.id}
              type="button"
              className={"add codegen-tab" + (l.id === language ? " active" : "")}
              onClick={() => setLanguage(l.id)}
            >
              {l.label}
            </button>
          ))}
          <button type="button" className="add" onClick={copy} disabled={!code}>
            {copied ? "Copied" : "Copy"}
          </button>
        </span>
      </h3>
      {busy && <p className="loading">Generating…</p>}
      {error && <p className="error">{error}</p>}
      {!busy && !error && <pre className="codegen-snippet">{code}</pre>}
    </div>
  );
}

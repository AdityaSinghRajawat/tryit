// Inline form to create a new stored secret. Per IMPL §8.3, the value
// goes panel→server only (never to the page).

import { useState } from "react";
import { ServerError, serverClient } from "../serverClient";
import type { SecretInfo, SecretType } from "../../shared/types";

interface Props {
  suggestedName: string;
  suggestedType: SecretType;
  suggestedHostHint?: string;
  onCreated: (info: SecretInfo) => void;
  onCancel: () => void;
}

export function SecretPrompt({
  suggestedName,
  suggestedType,
  suggestedHostHint,
  onCreated,
  onCancel,
}: Props) {
  const [name, setName] = useState(suggestedName);
  const [type, setType] = useState<SecretType>(suggestedType);
  const [hostHint, setHostHint] = useState(suggestedHostHint ?? "");
  const [value, setValue] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await serverClient.createSecret({
        name: name.trim(),
        type,
        hostHint: hostHint.trim() || undefined,
        value: type !== "basic" ? value : undefined,
        username: type === "basic" ? username : undefined,
        password: type === "basic" ? password : undefined,
      });
      onCreated({ name: name.trim(), type, hostHint: hostHint.trim() || undefined });
    } catch (err) {
      const msg = err instanceof ServerError ? `${err.code}: ${err.message}` : String(err);
      setError(msg);
    } finally {
      setBusy(false);
    }
  }

  const canSubmit =
    !!name.trim() &&
    (type === "basic" ? !!username && !!password : !!value);

  return (
    <form className="dialog secret-prompt" onSubmit={submit}>
      <h2>Add secret</h2>
      <p className="hint">
        Stored on this machine via the configured secrets backend. The value
        never leaves the panel→server channel.
      </p>

      <div className="row">
        <label>
          Name (UPPER_SNAKE_CASE)
          <input
            value={name}
            onChange={(e) => setName(e.target.value.toUpperCase())}
            spellCheck={false}
            autoFocus
          />
        </label>
      </div>

      <div className="row">
        <label>
          Type
          <select value={type} onChange={(e) => setType(e.target.value as SecretType)}>
            <option value="bearer">bearer</option>
            <option value="apiKey">apiKey</option>
            <option value="basic">basic</option>
          </select>
        </label>
      </div>

      <div className="row">
        <label>
          Host hint (optional — used for auto-matching)
          <input
            value={hostHint}
            onChange={(e) => setHostHint(e.target.value.trim())}
            placeholder="api.stripe.com"
            spellCheck={false}
          />
        </label>
      </div>

      {type !== "basic" && (
        <div className="row">
          <label>
            Value
            <input
              type="password"
              value={value}
              onChange={(e) => setValue(e.target.value)}
              spellCheck={false}
            />
          </label>
        </div>
      )}

      {type === "basic" && (
        <>
          <div className="row">
            <label>
              Username
              <input
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                spellCheck={false}
              />
            </label>
          </div>
          <div className="row">
            <label>
              Password
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                spellCheck={false}
              />
            </label>
          </div>
        </>
      )}

      <div className="send-row dialog-actions">
        <button type="button" className="add" onClick={onCancel} disabled={busy}>
          Cancel
        </button>
        <button type="submit" disabled={busy || !canSubmit}>
          {busy ? "Saving…" : "Save secret"}
        </button>
      </div>
      {error && <p className="error">{error}</p>}
    </form>
  );
}

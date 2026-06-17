// First-use consent gate per IMPL §8.4. Triggered when /execute returns a
// 200 with consentRequired:{secret, host}; on approve, POST /consent and
// let the caller retry.

import { useState } from "react";
import { ServerError, serverClient } from "../serverClient";

interface Props {
  secret: string;
  host: string;
  onGranted: () => void;
  onCancel: () => void;
}

export function ConsentDialog({ secret, host, onGranted, onCancel }: Props) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function grant() {
    setBusy(true);
    setError(null);
    try {
      await serverClient.grantConsent({ secret, host });
      onGranted();
    } catch (err) {
      setError(err instanceof ServerError ? `${err.code}: ${err.message}` : String(err));
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="dialog">
      <h2>Allow secret use?</h2>
      <p className="hint">
        Use <code>{secret}</code> against <code>{host}</code>?
        Tryit will remember this choice until you delete the secret or run{" "}
        <code>tryit reset-consent</code>.
      </p>
      <div className="send-row dialog-actions">
        <button type="button" className="add" onClick={onCancel} disabled={busy}>
          Cancel
        </button>
        <button type="button" onClick={grant} disabled={busy}>
          {busy ? "Granting…" : "Allow"}
        </button>
      </div>
      {error && <p className="error">{error}</p>}
    </div>
  );
}

import { useState } from "react";
import { ServerError, serverClient } from "../serverClient";

interface Props {
  onPaired: () => void;
}

export function PairDialog({ onPaired }: Props) {
  const [token, setToken] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await serverClient.pair(token.trim());
      onPaired();
    } catch (err) {
      const msg = err instanceof ServerError ? `${err.code}: ${err.message}` : String(err);
      setError(msg);
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="dialog">
      <h2>Pair with the tryit server</h2>
      <p className="hint">
        Run <code>make server-dev</code>, then paste the pairing token it
        printed to stdout. This binds this extension's origin to the server.
      </p>
      <form onSubmit={submit}>
        <label>
          Pairing token
          <input
            type="password"
            autoFocus
            value={token}
            onChange={(e) => setToken(e.target.value)}
            placeholder="tk_…"
            spellCheck={false}
          />
        </label>
        <button type="submit" disabled={busy || !token.trim()}>
          {busy ? "Pairing…" : "Pair"}
        </button>
      </form>
      {error && <p className="error">Pairing failed: {error}</p>}
    </div>
  );
}

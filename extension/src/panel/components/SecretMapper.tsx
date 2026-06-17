// Maps the spec's logical placeholders (e.g. STRIPE_API_KEY) to stored
// secrets. Auto-matches by hostHint per IMPL §8.3; the user can override
// or add a new secret inline. Output flows into ExecuteRequest.secretRefs.

import { useEffect, useMemo, useState } from "react";
import { ServerError, serverClient } from "../serverClient";
import { extractSecretRefs, hostOf } from "../../shared/secretRefs";
import type { RequestSpec, SecretInfo, SecretType } from "../../shared/types";
import { SecretPrompt } from "./SecretPrompt";

interface Props {
  spec: RequestSpec;
  onReady: (refs: Record<string, string>) => void;
  onCancel: () => void;
}

export function SecretMapper({ spec, onReady, onCancel }: Props) {
  const placeholders = useMemo(() => extractSecretRefs(spec), [spec]);
  const host = useMemo(() => hostOf(spec.baseUrl), [spec]);
  const authType: SecretType = spec.auth.type === "none" ? "bearer" : spec.auth.type;

  const [stored, setStored] = useState<SecretInfo[] | null>(null);
  const [refs, setRefs] = useState<Record<string, string>>({});
  const [creatingFor, setCreatingFor] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const r = await serverClient.listSecrets();
        setStored(r.secrets);
      } catch (err) {
        setError(err instanceof ServerError ? `${err.code}: ${err.message}` : String(err));
        setStored([]);
      }
    })();
  }, []);

  // Auto-suggest mapping once stored secrets arrive.
  useEffect(() => {
    if (!stored) return;
    setRefs((prev) => {
      const next = { ...prev };
      for (const ph of placeholders) {
        if (next[ph]) continue;
        const auto = autoMatch(stored, ph, host, authType);
        if (auto) next[ph] = auto.name;
      }
      return next;
    });
  }, [stored, placeholders, host, authType]);

  if (placeholders.length === 0) {
    // No placeholders — caller can skip straight to send.
    return null;
  }

  if (stored === null) {
    return <p className="loading">Loading secrets…</p>;
  }

  if (creatingFor !== null) {
    return (
      <SecretPrompt
        suggestedName={creatingFor}
        suggestedType={authType}
        suggestedHostHint={host || undefined}
        onCreated={(info) => {
          setStored((prev) => (prev ? [...prev, info] : [info]));
          setRefs((prev) => ({ ...prev, [creatingFor]: info.name }));
          setCreatingFor(null);
        }}
        onCancel={() => setCreatingFor(null)}
      />
    );
  }

  const allMapped = placeholders.every((ph) => !!refs[ph]);

  return (
    <div className="dialog">
      <h2>Map secrets</h2>
      <p className="hint">
        Each placeholder picks a secret stored on this machine. Auto-matched by
        host where possible; override or add a new one.
      </p>

      {placeholders.map((ph) => {
        const candidates = candidatesFor(stored, ph, host, authType);
        const value = refs[ph] ?? "";
        return (
          <div key={ph} className="row kv-row">
            <span className="kv-name"><code>{ph}</code></span>
            <select
              className="kv-value"
              value={value}
              onChange={(e) => setRefs((prev) => ({ ...prev, [ph]: e.target.value }))}
            >
              <option value="" disabled>
                — pick a secret —
              </option>
              {candidates.map((c) => (
                <option key={c.name} value={c.name}>
                  {c.name}
                  {c.hostHint ? ` (${c.hostHint})` : ""}
                </option>
              ))}
              {stored
                .filter((s) => !candidates.find((c) => c.name === s.name))
                .map((s) => (
                  <option key={s.name} value={s.name}>
                    {s.name}
                    {s.hostHint ? ` (${s.hostHint})` : ""}
                    {" · other"}
                  </option>
                ))}
            </select>
            <button
              type="button"
              className="add"
              onClick={() => setCreatingFor(ph)}
            >
              New…
            </button>
          </div>
        );
      })}

      <div className="send-row dialog-actions">
        <button type="button" className="add" onClick={onCancel}>
          Back
        </button>
        <button
          type="button"
          disabled={!allMapped}
          onClick={() => onReady(refs)}
        >
          Continue
        </button>
      </div>
      {error && <p className="error">{error}</p>}
    </div>
  );
}

// candidatesFor ranks stored secrets for a given placeholder: a hostHint
// match beats everything; an auth-type match beats nothing; otherwise the
// secret falls into the "other" bucket the caller appends.
function candidatesFor(
  stored: SecretInfo[],
  _placeholder: string,
  host: string,
  authType: SecretType,
): SecretInfo[] {
  const byHostAndType: SecretInfo[] = [];
  const byHost: SecretInfo[] = [];
  const byType: SecretInfo[] = [];

  for (const s of stored) {
    const hostMatch = !!host && !!s.hostHint && host.endsWith(s.hostHint);
    const typeMatch = s.type === authType;
    if (hostMatch && typeMatch) byHostAndType.push(s);
    else if (hostMatch) byHost.push(s);
    else if (typeMatch) byType.push(s);
  }
  return [...byHostAndType, ...byHost, ...byType];
}

function autoMatch(
  stored: SecretInfo[],
  placeholder: string,
  host: string,
  authType: SecretType,
): SecretInfo | undefined {
  const ranked = candidatesFor(stored, placeholder, host, authType);
  return ranked[0];
}

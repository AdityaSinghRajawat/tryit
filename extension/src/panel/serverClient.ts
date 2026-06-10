// D13: the panel is the SOLE caller of the Go server. It owns the pair
// token, and attaches it to every protected route.

import {
  SERVER_BASE,
  STORAGE_KEY_PAIR_TOKEN,
} from "../shared/constants";
import type {
  ExecuteRequest,
  ExecuteResponse,
  HealthResponse,
  PairRequest,
  PairResponse,
  ErrorEnvelope,
} from "../shared/types";

export class ServerError extends Error {
  code: string;
  status: number;
  constructor(code: string, message: string, status: number) {
    super(message);
    this.code = code;
    this.status = status;
  }
}

async function getToken(): Promise<string | undefined> {
  const r = await chrome.storage.local.get(STORAGE_KEY_PAIR_TOKEN);
  return r[STORAGE_KEY_PAIR_TOKEN] as string | undefined;
}

async function call<T>(
  path: string,
  init: RequestInit,
  needsToken: boolean
): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set("Content-Type", "application/json");
  if (needsToken) {
    const tok = await getToken();
    if (!tok) {
      throw new ServerError("not_paired", "extension is not paired yet", 409);
    }
    headers.set("Authorization", "Bearer " + tok);
  }
  let resp: Response;
  try {
    resp = await fetch(SERVER_BASE + path, { ...init, headers });
  } catch (e) {
    throw new ServerError(
      "target_unreachable",
      "could not reach the tryit server at " + SERVER_BASE,
      0
    );
  }
  const text = await resp.text();
  if (!resp.ok) {
    let env: ErrorEnvelope | undefined;
    try {
      env = text ? JSON.parse(text) : undefined;
    } catch {
      /* not JSON */
    }
    const code = env?.error?.code ?? "internal";
    const msg = env?.error?.message ?? "request failed";
    throw new ServerError(code, msg, resp.status);
  }
  return (text ? JSON.parse(text) : {}) as T;
}

export const serverClient = {
  health(): Promise<HealthResponse> {
    return call<HealthResponse>("/health", { method: "GET" }, false);
  },

  async pair(token: string): Promise<PairResponse> {
    const req: PairRequest = { token };
    const r = await call<PairResponse>(
      "/pair",
      { method: "POST", body: JSON.stringify(req) },
      false
    );
    // The panel learns the token by storing it locally; only then are protected calls possible.
    await chrome.storage.local.set({ [STORAGE_KEY_PAIR_TOKEN]: token });
    return r;
  },

  async forgetToken(): Promise<void> {
    await chrome.storage.local.remove(STORAGE_KEY_PAIR_TOKEN);
  },

  async hasToken(): Promise<boolean> {
    return (await getToken()) !== undefined;
  },

  execute(req: ExecuteRequest): Promise<ExecuteResponse> {
    return call<ExecuteResponse>(
      "/execute",
      { method: "POST", body: JSON.stringify(req) },
      true
    );
  },
};

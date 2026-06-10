# DECISIONS

Log of decisions made during implementation where both SPEC and IMPL leave a small gap. Each entry mirrors a `// DECISION:` comment in the code.

## Phase 1

### D-P1-1 — Module path
`github.com/AdityaSinghRajawat/tryit/server` (matches the git remote `github.com:AdityaSinghRajawat/tryit`). Easy to change later via `go mod edit -module ...` + sed across imports.

### D-P1-2 — Minimal `/pair` endpoint in Phase 1
The Phase 4 file list in §4 lists `pairHandler` under Phase 3, but the Phase 1 **Verify** step explicitly requires "pair (paste token)" and §8.2 says protected routes 409 if not paired. So Phase 1 ships a minimal `pairService` + `pairHandler` that implement §8.1 end-to-end (token compare → bind origin). Phase 3 adds rotation, reset commands, and CLI ergonomics.

### D-P1-3 — `pairStore` Phase 1 backend = file, not keychain
§8.1 says keychain entry `tryit-pair-token`, but Phase 1 is "no-keychain". Use `~/.tryit/pair.json` (0600) for Phase 1; swap to keychain in Phase 2 alongside the secrets backend swap. Same `pair.Store` port, one-line change in `main.go`.

### D-P1-4 — Phase 1 has no `/parse` endpoint
Phase 1 verify is "click Try it on a Swagger UI page → live response". The content script harvests the OpenAPI operation object as `structuredHint` and the panel converts it into a `RequestSpec` client-side (Phase 1 only). Phase 2 introduces `/parse` and moves that logic server-side. Rationale: keeps Phase 1 truly no-AI / no-cascade and proves the cross-origin + secret-stays-server-side pipeline.

### D-P1-5 — Phase 1 secret source = `TRYIT_SECRET_<NAME>` env vars
`storage.envStore` looks up `TRYIT_SECRET_<UPPER_SNAKE_NAME>`; missing → `secret_not_found`. Phase 2 swaps to keychain. The Phase 1 dev sets `TRYIT_SECRET_PETSTORE_KEY=special-key` in the server process env.

### D-P1-6 — Consent gate deferred to Phase 2
§8.4's per-(secret, host) consent flow ships in Phase 2 with the secrets manager. Phase 1 has no `consentStore`; `executeService` skips the consent check. The Phase 1 secret comes from env, so first-use risk is bounded.

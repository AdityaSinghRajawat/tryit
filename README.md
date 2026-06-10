# tryit

Make any API docs page executable inline.

> Phase 1 (current): vertical slice — Swagger UI detection only, no AI, secret from env. See `DECISIONS.md`, `tryit-SPEC.txt`, `tryit-IMPL.txt`.

## Quick start (development)

Requires Go 1.22+ and Node 18+.

```bash
# Terminal 1 — Go server (prints the pairing token on first start)
make server-dev

# Terminal 2 — extension watch build
make ext-install
make ext-dev
```

Then in Chromium:

1. `chrome://extensions` → enable Developer Mode → **Load unpacked** → select `extension/dist`.
2. Open the extension's side panel and paste the pairing token from the server's stdout.
3. Visit a Swagger UI page (e.g. https://petstore.swagger.io/) → click a **Try it** button injected next to an endpoint → see the live response in the panel.

## Layout

```
server/             Go single binary (api/, cmd/, internal/...)
extension/          TS / Vite / CRXJS / React
fixtures/           Golden test inputs/expected (Phase 2+)
```

License: Apache-2.0 (see `LICENSE`, added in Phase 3).

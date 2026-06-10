import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { crx } from "@crxjs/vite-plugin";
import manifest from "./src/manifest.json" with { type: "json" };

export default defineConfig({
  plugins: [react(), crx({ manifest: manifest as any })],
  build: {
    outDir: "dist",
    emptyOutDir: true,
    target: "esnext",
  },
});

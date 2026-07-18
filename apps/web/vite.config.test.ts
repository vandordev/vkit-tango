import { expect, test } from "bun:test";

test("configures TanStack Start with Bun-targeted Nitro", async () => {
  const viteConfig = await Bun.file(new URL("./vite.config.ts", import.meta.url)).text();

  expect(viteConfig).toContain('from "@tanstack/react-start/plugin/vite"');
  expect(viteConfig).toContain('from "@tanstack/router-plugin/vite"');
  expect(viteConfig).toContain('from "nitro/vite"');
  expect(viteConfig).toContain('nitro({ preset: "bun" })');
  expect(viteConfig).toContain("tanstackStart()");
  expect(viteConfig).toContain('tanstackRouter({ target: "react", autoCodeSplitting: true })');
});

test("uses the YAML wrapper without a Next.js command", async () => {
  const { scripts } = await Bun.file(new URL("./package.json", import.meta.url)).json();

  expect(scripts.dev).toContain("--env-file=../../.env");
  expect(scripts.dev).toContain("--modules base,web,api,storage");
  expect(scripts.dev).toContain("vite --port 4100");
  expect(scripts.build).toContain("vite build");
  expect(scripts.start).toContain(".output/server/index.mjs");
  expect(scripts.dev).not.toContain("next");
  expect(scripts.build).not.toContain("next");
});

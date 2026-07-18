import { expect, test } from "bun:test";

test("documents TanStack Start as the web runtime", async () => {
  const readme = await Bun.file(new URL("../../README.md", import.meta.url)).text();
  const webRules = await Bun.file(new URL("../../.agent/web.md", import.meta.url)).text();

  expect(readme).toContain("TanStack Start");
  expect(readme).not.toContain("Next.js for the web experience");
  expect(webRules).toContain("TanStack Start");
  expect(webRules).toContain("src/server.ts");
  expect(webRules).not.toContain("App Router");
});

import { readFileSync } from "node:fs";
import { join } from "node:path";
import { expect, test } from "bun:test";

const root = join(import.meta.dir, "..");

test("removes unused Next-specific lint support", () => {
  const eslintPackage = Bun.file(join(root, "packages/eslint-config/package.json"));

  return eslintPackage.json().then(({ devDependencies, exports }) => {
    expect(devDependencies["@next/eslint-plugin-next"]).toBeUndefined();
    expect(exports["./next-js"]).toBeUndefined();
  });
});

test("documents Vite public configuration in security guidance", () => {
  const security = readFileSync(join(root, "SECURITY.md"), "utf8");

  expect(security).toContain("VITE_");
  expect(security).not.toContain("NEXT_PUBLIC_");
});

test("uses the Elysia health endpoint for the Compose web healthcheck", () => {
  const compose = readFileSync(join(root, "docker-compose.yml"), "utf8");

  expect(compose).toContain("http://127.0.0.1:4100/health");
});

import { readFileSync } from "node:fs";
import { join } from "node:path";
import { expect, test } from "bun:test";

const root = join(import.meta.dir, "..");

test("uses vkit-fast for public and Compose identities", () => {
  const previousName = "vkit-" + "rapid";
  const readme = readFileSync(join(root, "README.md"), "utf8");
  const compose = readFileSync(join(root, "docker-compose.yml"), "utf8");

  expect(readme).toContain("# vkit-fast");
  expect(readme).not.toContain(previousName);
  expect(compose).toMatch(/^name: vkit-fast$/m);
});

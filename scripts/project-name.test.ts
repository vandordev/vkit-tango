import { readFileSync } from "node:fs";
import { join } from "node:path";
import { expect, test } from "bun:test";

const root = join(import.meta.dir, "..");

test("uses vkit-tango for public and Compose identities", () => {
  const previousName = "vkit-" + "rapid";
  const readme = readFileSync(join(root, "README.md"), "utf8");
  const compose = readFileSync(join(root, "docker-compose.yml"), "utf8");
  const appConfig = readFileSync(join(root, "config/app.yaml"), "utf8");

  expect(readme).toContain("# vkit-tango");
  expect(readme).not.toContain(previousName);
  expect(compose).toMatch(/^name: vkit-tango$/m);
  expect(appConfig).toContain("name: vkit-tango");
});

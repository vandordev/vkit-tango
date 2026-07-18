import { expect, test } from "bun:test";

test("documents shadcn ui as the default and alternatives as project choices", async () => {
  const [readme, agentRules, uiRules] = await Promise.all([
    Bun.file(new URL("../../README.md", import.meta.url)).text(),
    Bun.file(new URL("../../AGENTS.md", import.meta.url)).text(),
    Bun.file(new URL("../../.agent/ui.md", import.meta.url)).text(),
  ]);

  for (const document of [readme, agentRules, uiRules]) {
    expect(document).toContain("shadcn/ui");
    expect(document).toContain("one primary UI system");
    expect(document).toContain("Mantine");
    expect(document).toContain("MUI");
  }

  expect(readme).toContain("TanStack Start");
  expect(readme).not.toContain("Mantine by default");
  expect(agentRules).not.toContain("Mantine is the current default");
});

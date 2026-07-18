import { expect, test } from "bun:test";

test("provides English not-found and error recovery actions", async () => {
  const rootRoute = await Bun.file(new URL("./__root.tsx", import.meta.url)).text();

  expect(rootRoute).toContain("Page not found");
  expect(rootRoute).toContain("Back to home");
  expect(rootRoute).toContain("Something went wrong");
  expect(rootRoute).toContain("Try again");
  expect(rootRoute).toContain("onClick={reset}");
  expect(rootRoute).toContain("import.meta.env.DEV");
  expect(rootRoute).toContain('from "@/components/ui/button"');
  expect(rootRoute).not.toContain("@mantine/");
  expect(rootRoute).not.toContain("MantineProvider");
  expect(rootRoute).not.toContain("ColorSchemeScript");
  expect(rootRoute).not.toContain("Notifications");
});

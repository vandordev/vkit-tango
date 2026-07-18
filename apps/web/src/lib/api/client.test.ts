import { expect, test } from "bun:test";

test("builds an Eden client at the same-origin API base path", async () => {
  const { api, createApiClient } = await import("./client");

  expect(api).toBeDefined();
  expect(createApiClient("/api")).toBeDefined();
});

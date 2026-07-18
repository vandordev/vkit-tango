import { expect, test } from "bun:test";

test("builds a generated fetch client at the same-origin API base path", async () => {
  const { api, createApiClient } = await import("./client");

  expect(api).toBeDefined();
  expect(createApiClient("").getConfig().baseUrl).toBe("");
});

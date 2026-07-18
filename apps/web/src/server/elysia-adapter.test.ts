import { expect, test } from "bun:test";

import { resolvedConfigEnvironment } from "../../../../packages/config/src/run";

Object.assign(
  process.env,
  resolvedConfigEnvironment(["base", "api", "web"], {
    DATABASE_URL: "postgresql://db",
    NODE_ENV: "test",
  }),
);

test("forwards an API request to embedded Elysia", async () => {
  const { forwardApiRequest } = await import("./elysia-adapter");
  const response = await forwardApiRequest(new Request("http://localhost:4100/api/status"));

  expect(response.status).toBe(200);
  expect(await response.json()).toEqual({ success: true, data: { status: "ok" } });
  expect(response.headers.get("x-request-id")).toBeString();
});

test("forwards health requests to embedded Elysia", async () => {
  const { forwardHealthRequest } = await import("./elysia-adapter");
  const response = await forwardHealthRequest(new Request("http://localhost:4100/health"));

  expect(response.status).toBe(200);
  expect((await response.json()).data.status).toBe("healthy");
});

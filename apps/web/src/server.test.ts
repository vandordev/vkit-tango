import { expect, test } from "bun:test";

import { resolvedConfigEnvironment } from "../../../packages/config/src/run";

Object.assign(
  process.env,
  resolvedConfigEnvironment(["base", "api", "web"], {
    DATABASE_URL: "postgresql://db",
    NODE_ENV: "test",
  }),
);

test("forwards embedded API and health requests before the Start handler", async () => {
  const { default: server } = await import("./server");

  const [health, status] = await Promise.all([
    server.fetch(new Request("http://localhost:4100/health")),
    server.fetch(new Request("http://localhost:4100/api/status")),
  ]);

  expect(health.status).toBe(200);
  expect(status.status).toBe(200);
  expect((await health.json()).data.status).toBe("healthy");
  expect(await status.json()).toEqual({ success: true, data: { status: "ok" } });
});

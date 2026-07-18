import { expect, test } from "bun:test";

import { createRealtimeConfig } from "./realtime";
import { loadConfig } from "./loader";

const configDirectory = new URL("../../../config", import.meta.url).pathname;

test("creates a scoped realtime runtime configuration", () => {
  expect(
    createRealtimeConfig(
      loadConfig({
        configDirectory,
        modules: ["realtime"],
        environment: {
          REALTIME_TICKET_SECRET: "ticket-secret",
          REALTIME_INTERNAL_API_KEY: "publisher-key",
        },
      }) as Record<string, string | undefined>,
    ),
  ).toMatchObject({ port: 4102 });
});

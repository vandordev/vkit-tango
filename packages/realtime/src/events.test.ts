import { expect, test } from "bun:test";

import { realtimeEventSchema, roomsForEvent } from "./events";

test("routes a resource event to resource and workspace rooms", () => {
  const event = realtimeEventSchema.parse({
    type: "resource.updated.v1",
    event_id: crypto.randomUUID(),
    occurred_at: new Date().toISOString(),
    resource_id: "r1",
    workspace_id: "w1",
  });

  expect(roomsForEvent(event)).toEqual(["resource:r1", "workspace:w1"]);
});

test("rejects the unversioned resource event", () => {
  expect(() =>
    realtimeEventSchema.parse({
      type: "resource.updated",
      event_id: crypto.randomUUID(),
      occurred_at: new Date().toISOString(),
      resource_id: "r1",
      workspace_id: "w1",
    }),
  ).toThrow();
});

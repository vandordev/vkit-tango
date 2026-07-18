import { z } from "zod";

export const realtimeEventSchema = z.object({
  type: z.literal("resource.updated.v1"),
  event_id: z.string().uuid(),
  occurred_at: z.string().datetime(),
  resource_id: z.string().min(1),
  workspace_id: z.string().min(1),
});

export type RealtimeEvent = z.infer<typeof realtimeEventSchema>;

export function roomsForEvent(event: RealtimeEvent): string[] {
  return [`resource:${event.resource_id}`, `workspace:${event.workspace_id}`];
}

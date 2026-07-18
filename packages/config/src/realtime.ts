import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

import { commonServer } from "./common";

const realtimeServer = {
  ...commonServer,
  REALTIME_PORT: z.coerce.number().int().positive().default(4102),
  REALTIME_TICKET_SECRET: z.string().min(1),
  REALTIME_PUBLISH_API_KEY: z.string().min(1),
} as const;

export function createRealtimeConfig(runtimeEnv: Record<string, unknown>) {
  const realtime = runtimeEnv.realtime;
  const resolvedEnv =
    realtime && typeof realtime === "object" && !Array.isArray(realtime)
      ? {
          NODE_ENV: typeof runtimeEnv.NODE_ENV === "string" ? runtimeEnv.NODE_ENV : undefined,
          DATABASE_URL: typeof runtimeEnv.DATABASE_URL === "string" ? runtimeEnv.DATABASE_URL : undefined,
          REALTIME_PORT: (realtime as Record<string, unknown>).port,
          REALTIME_TICKET_SECRET: (realtime as Record<string, unknown>).ticket_secret,
          REALTIME_PUBLISH_API_KEY: (realtime as Record<string, unknown>).internal_api_key,
        }
      : runtimeEnv;

  const parsed = createEnv({
    server: realtimeServer,
    runtimeEnv: resolvedEnv as Record<string, string | undefined>,
    isServer: true,
    emptyStringAsUndefined: true,
  });

  return {
    ...parsed,
    port: parsed.REALTIME_PORT,
  };
}

export type RealtimeConfig = ReturnType<typeof createRealtimeConfig>;

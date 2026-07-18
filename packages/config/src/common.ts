import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

export const commonServer = {
  NODE_ENV: z
    .enum(["development", "staging", "production", "test"])
    .default("development"),
  LOG_LEVEL: z.string().default("info"),
} as const;

export function createCommonConfig(
  runtimeEnv: Record<string, string | undefined>,
) {
  const parsed = createEnv({
    server: commonServer,
    runtimeEnv,
    isServer: true,
    emptyStringAsUndefined: true,
  });

  return parsed;
}

export type CommonConfig = ReturnType<typeof createCommonConfig>;

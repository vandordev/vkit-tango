import { treaty } from "@elysia/eden";

import type { App } from "@repo/api";

type ApiClient = ReturnType<typeof treaty<App>>["api"];

export function createApiClient(baseUrl: string): ApiClient {
  return treaty<App>(baseUrl).api;
}

export const api: ApiClient = createApiClient("/api");

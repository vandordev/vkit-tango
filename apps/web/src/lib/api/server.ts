import { createApiClient } from "./client";

export const api = createApiClient(process.env.API_BASE_URL ?? "http://localhost:4101");

import { createClient } from "./generated/client";
import { client } from "./generated/client.gen";

export function createApiClient(baseUrl: string) {
  return createClient({ baseUrl });
}

client.setConfig({ baseUrl: "" });

export { client as api };

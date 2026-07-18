import { createFileRoute } from "@tanstack/react-router";

import { forwardApiRequest } from "@/server/elysia-adapter";

const handler = ({ request }: { request: Request }) => forwardApiRequest(request);

export const Route = createFileRoute("/api/$")({
  server: {
    handlers: {
      GET: handler,
      POST: handler,
      PUT: handler,
      PATCH: handler,
      DELETE: handler,
      OPTIONS: handler,
    },
  },
});

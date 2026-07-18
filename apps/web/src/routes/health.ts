import { createFileRoute } from "@tanstack/react-router";

import { forwardHealthRequest } from "@/server/elysia-adapter";

export const Route = createFileRoute("/health")({
  server: {
    handlers: {
      GET: ({ request }) => forwardHealthRequest(request),
    },
  },
});

import handler, { createServerEntry } from "@tanstack/react-start/server-entry";

import { forwardApiRequest, forwardHealthRequest } from "./server/elysia-adapter";

export default createServerEntry({
  fetch(request) {
    const pathname = new URL(request.url).pathname;

    if (pathname === "/health") {
      return forwardHealthRequest(request);
    }

    if (pathname === "/api" || pathname.startsWith("/api/")) {
      return forwardApiRequest(request);
    }

    return handler.fetch(request);
  },
});

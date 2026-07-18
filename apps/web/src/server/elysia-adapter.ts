import { app } from "@repo/api";

export function forwardApiRequest(request: Request): Response | Promise<Response> {
  return app.fetch(request);
}

export function forwardHealthRequest(request: Request): Response | Promise<Response> {
  return app.fetch(request);
}

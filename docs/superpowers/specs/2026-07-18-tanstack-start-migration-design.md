# TanStack Start Migration Design

## Goal

Replace the `apps/web` Next.js runtime with TanStack Start while preserving the boilerplate's existing product boundaries: Elysia remains the only HTTP API, Eden remains the only frontend API client, and the default web process continues to embed the API at the same public origin.

The migration deliberately does not add a product domain, authentication, a second API transport, Tailwind, or shadcn/ui. Mantine remains the web UI baseline.

## Selected approach

`apps/web` becomes a Vite application using TanStack Start, TanStack Router, and Nitro with Bun as the production preset. TanStack Start's custom server entry owns only the framework adapter layer:

```text
Browser -> Eden -> /api/* -> TanStack Start server entry -> Elysia app.fetch -> API route
Browser -> GET /health -> TanStack Start server entry -> Elysia app.fetch -> health route
```

`src/server.ts` forwards matching incoming `Request` objects unchanged to `@repo/api`'s exported `app.fetch` handler and returns its `Response` unchanged. It supports the same HTTP methods as the former adapter for `/api/*` and forwards `/health`; all other requests go to TanStack Start's default handler.

This preserves the public origin and endpoint paths (`/api/status`, `/api/docs`, `/api/openapi.json`, and `/health`), avoids a reverse proxy and a second default process, and keeps standalone Elysia at port 4101 available for independent deployments.

## Web application structure

The former `app/` directory is replaced with TanStack Router source under `apps/web/src/`:

- `router.tsx` creates and registers the TanStack Router instance from generated routes. It enables scroll restoration and intent preloading.
- `routes/__root.tsx` owns the HTML document, metadata, Mantine providers, the React Query provider, notifications, global styles, and the TanStack `HeadContent` and `Scripts` components.
- `routes/index.tsx` is the public landing route and `routes/dashboard.tsx` is the dashboard route. They preserve the current content and paths.
- `server.ts` is the Elysia adapter in TanStack Start's custom server-entry extension point. It intercepts `/api/*` and `/health`; it does not define a UI route.
- `lib/api/client.ts` remains the browser-only Eden facade. It calls the same-origin `/api` base path and imports only `type App` from `@repo/api`.
- `lib/api/server.ts` remains server-only. It uses `treaty(app).api` for server-rendered loaders or server functions, rather than sending an unnecessary loopback request.

No web component imports Prisma, an application usecase, or a non-type API implementation. Server-only Elysia imports stay inside route handlers and server-only modules so they are excluded from browser bundles.

The generated `routeTree.gen.ts` is a build artifact produced by the TanStack Router plugin. It is not edited by hand and is included in source control only if the repository's generated-file convention requires it.

## Dependencies, build, and runtime

`apps/web/package.json` removes `next`, `@t3-oss/env-nextjs`, `server-only`, and Next-specific lint/typecheck support. It adds pinned compatible versions of Vite, TanStack Start, TanStack Router, the Router Vite plugin, and Nitro. The migration uses the installed `../my-tanstack-start` only as a structural reference; it does not copy its Tailwind, shadcn, Prisma, demo, or floating `latest` dependency choices.

`vite.config.ts` composes TanStack Start, the Router plugin, React, and Nitro. Nitro is configured with the Bun preset so the production web process runs on Bun, matching the rest of the repository. Workspace API packages required by the embedded adapter are explicitly handled by the SSR build configuration so `@repo/api` is not accidentally emitted as a browser dependency.

The web scripts retain their existing names and ports:

- `dev` starts Vite on port 4100 through the YAML configuration wrapper.
- `build` creates the TanStack Start/Nitro production output through the same wrapper.
- `start` runs the generated Bun server on port 4100 through the wrapper.
- `check-types` runs TypeScript and route generation/type validation; `lint` uses the repository's non-Next ESLint baseline.
- `clean` removes Vite, Nitro, and router-generated build artifacts rather than `.next`.

Root package scripts, Turbo task discovery, Taskfile commands, the web Dockerfile, and Compose continue to expose the same `task dev`, `task start`, `task build`, `task web:health`, and Compose workflows. The web image builds and starts the generated Bun-compatible TanStack Start output rather than a Next standalone server.

## Configuration and security

The YAML-first configuration wrapper remains the sole mechanism that resolves `config/` modules before a web build or server starts. The embedded web process continues to select `base,web,api,storage`, which provides server-side Elysia with `DATABASE_URL` and OpenAPI configuration without exposing those values to browser code.

The client Eden facade uses the fixed same-origin `/api` base path. Therefore `NEXT_PUBLIC_APP_URL` is removed from `config/web.yaml`, `.env.example`, Compose, and web configuration tests; a replacement `VITE_*` public URL is not needed for the baseline. If a future browser-visible setting is required, it must use Vite's `VITE_` prefix, be validated in an `@t3-oss/env-core` schema compatible with Vite, and be documented as public.

The API's `CORS_ORIGIN` continues to default to `http://localhost:4100`, and `OPENAPI_SERVER_URL` continues to resolve to port 4100 in the embedded web configuration and port 4101 for the standalone API configuration. Elysia remains responsible for API validation, request logging, CORS, OpenAPI, and error envelopes; TanStack Start does not duplicate those concerns.

## Documentation and ownership rules

`README.md` is rewritten where it names the former framework, its route handlers, old environment semantics, or `.next` output. It documents TanStack Start, TanStack Router, Vite/Nitro, the embedded server-entry adapter, and the same command surface.

`.agent/web.md`, `.agent/architecture.md`, and `.agent/config.md` are updated to replace Next-specific rules with TanStack Start equivalents. The new rules state that:

- Elysia is embedded only through the TanStack Start custom server entry.
- Browser code reaches Elysia only through Eden at same-origin `/api`.
- Server data access uses a dedicated server-only Eden facade.
- Browser-visible configuration uses only explicit `VITE_*` values; database and API secrets remain server-only.

The API, database, application, worker, scheduler, queue, storage, and realtime ownership documents remain unchanged except for wording that directly describes the old adapter.

## Testing and acceptance criteria

Focused tests are written before each production change. They verify:

- the server entry forwards API requests to Elysia and preserves the response status, body, and headers;
- the server entry returns Elysia's health response;
- the browser Eden facade remains typed and targets same-origin `/api`;
- the server Eden facade remains isolated from browser entry points;
- web configuration has no public database value and no dependency on `NEXT_PUBLIC_APP_URL`;
- the Vite configuration selects the Bun/Nitro runtime and route generation completes;
- the web Docker/Compose configuration builds and starts the new output without legacy-framework variables or commands.

The migration is accepted when `task test:web`, `task check-types:web`, and `task build:web` pass; `/health` and `/api/status` respond successfully from the web process on port 4100; and the repository-wide `task quality` and `task build` pass. The standalone `task dev:standalone-api` and its health endpoint remain operational.

## Non-goals

- Moving Elysia business routes into TanStack Start server functions or route handlers.
- Replacing Eden with TanStack Query fetchers, tRPC, or direct page-level `fetch` calls.
- Removing the standalone Elysia server entrypoint.
- Changing the database, queue, worker, scheduler, storage, or realtime implementation.
- Changing the UI system from Mantine to Tailwind or shadcn/ui.
- Supporting a deployment target other than the existing Bun-oriented container baseline in this migration.

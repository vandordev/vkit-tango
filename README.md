# vkit-tango

**Tango** stands for **TanStack Start + Go**—a monorepo for web, APIs, workers,
migrations, and realtime, all built on PostgreSQL.

A data-driven monorepo for TanStack Start, Go, PostgreSQL, and Socket.IO.

The web UI uses shadcn/ui as its one primary UI system. Choose Mantine or MUI
only deliberately for a specific project.

## Architecture

- `apps/web`: TanStack Start. The browser calls the Go API through `/api/*`
  and uses Hey API-generated clients and TanStack Query hooks.
- `apps/api`: Go HTTP API built with Huma. Every business endpoint uses
  `/api/v1/*`; `/health`, `/health/ready`, `/api/openapi.json`, and
  `/api/docs` are process-level endpoints.
- `apps/worker`: Go River worker. All background processing and River
  schedules run in Go.
- `apps/scheduler`: Go Fx periodic-enqueue process; it has no executable
  workers.
- `apps/migrate`: Single Goose and River migration process.
- `database/schema`: Ent schema source; generated client in
  `internal/platform/db`.
- `database/migrations`: Goose migrations.
- `internal/usecase`: The single source of business-mutation rules. Query
  handlers may read Ent directly; mutations must use a use case.
- `apps/realtime` and `packages/realtime`: TypeScript Socket.IO runtime.
  Go publishes events to the private `/internal/events` endpoint; its contract
  is `contracts/asyncapi/realtime.v1.yaml`.

PostgreSQL is the only queue backend. Ent mutations and River enqueueing must
share the same SQL transaction.

## Configuration

Snake_case YAML configuration lives in `config/`. Go explicitly loads the
modules it needs (`database`, `http_api`, `worker`, and `realtime`), while
TypeScript loads only `web` or `realtime`. Secrets are supplied exclusively
through environment interpolation.

## Commands

```bash
task install
cp .env.example .env
task migrate
task dev
task dev -- api web realtime
task dev:api
task dev:worker
task dev:scheduler
task dev:web
task dev:realtime
task api:client:generate
task quality
task build
```

`task dev` runs API, worker, scheduler, and web together. Select services with
arguments after `--`, for example `task dev -- api web realtime`; realtime is
not part of the default set. Run `task migrate` separately before starting the
long-lived services.

Generate Go surfaces only through Taskfile: `task add:usecase name=...`,
`task add:http-handler name=... method=PUT path=/api/v1/...`, `task add:job
name=...`, and `task add:scheduler name=...`. Refresh committed Fx registries
with `task sync:usecase`, `task sync:http`, `task sync:worker`,
`task sync:scheduler`, or `task sync`. The umbrella `task sync` also refreshes
the OpenAPI document and Hey API client. Only use cases require a paired test.

Use `docker compose up --build` to run PostgreSQL, migrations, and the Go API. Add `--profile jobs` for the worker.

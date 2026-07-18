# vkit-tango

**Tango** stands for **TanStack Start + Go**—a monorepo for web, APIs, workers,
migrations, and realtime, all built on PostgreSQL.

A data-driven monorepo for TanStack Start, Go, PostgreSQL, and Socket.IO.

The web UI uses shadcn/ui as its one primary UI system. Choose Mantine or MUI
only deliberately for a specific project.

## Architecture

- `apps/web`: TanStack Start. The browser calls the Go API through `/api/*`
  and uses Hey API-generated clients and TanStack Query hooks.
- `apps/api`: Uber Fx Go HTTP composition root built with Chi, Huma, and
  `humachi`. Every business endpoint uses
  `/api/v1/*`; `/health`, `/health/ready`, `/api/openapi.json`, and
  `/api/docs` are process-level endpoints.
- `apps/worker`: Uber Fx Go River worker composition root. It executes typed
  jobs that call commands.
- `apps/scheduler`: Uber Fx periodic-enqueue composition root; it has no
  executable workers and only enqueues typed jobs.
- `apps/migrate`: Single Goose and River migration process.
- `database/schema`: Ent schema source; generated client in
  `internal/platform/db`.
- `database/migrations`: Goose migrations.
- `internal/usecase`: The single source of business-mutation rules. Query
  handlers may read Ent directly; mutations must use a use case.
- `internal/contract`: Shared command, HTTP-handler, River-job, and scheduler
  registration boundaries. Generated Fx registries under `internal/generated/fx`
  wire valid implementations deterministically.
- `apps/realtime` and `packages/realtime`: TypeScript Socket.IO runtime.
  Go publishes events to the private `/internal/events` endpoint; its contract
  is `contracts/asyncapi/realtime.v1.yaml`.

PostgreSQL is the only queue backend. Ent mutations and River enqueueing must
share the same SQL transaction.

## Configuration

Snake_case YAML configuration lives in `config/`. Go explicitly loads the
modules it needs (`database`, `http_api`, `worker`, `scheduler`, and `realtime`), while
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

`vx` is an internal implementation detail of `task add:*`; developers and CI
use Taskfile commands only. Never edit `internal/generated/fx` by hand. Run a
focused test before each change, then run `task sync` after generator or
registry changes. Shared infrastructure or runtime changes additionally require
`task quality` and `task build`.

Use `docker compose up --build` to run PostgreSQL, migrations, and the Go API. Add `--profile jobs` for the worker.

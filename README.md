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

## Prerequisites

Install these tools and ensure their binary directories are on `PATH` before
running the project:

| Tool | Required version | Purpose |
| --- | --- | --- |
| [Go](https://go.dev/dl/) | `1.25.7` | Go runtimes, generators, and the project-managed Air tool. |
| [Bun](https://bun.com/docs/installation) | `1.3.14` | Web/realtime runtime, workspace dependencies, tests, and generators. |
| [Task](https://taskfile.dev/docs/installation) | v3 | The only developer and CI command interface. |
| [RTK](https://github.com/rtk-ai/rtk) | latest | Command proxy used by the Taskfile and AI-agent workflow. |
| [Docker Engine with Docker Compose](https://docs.docker.com/compose/install/) | Compose v2 | Local PostgreSQL and containerized runtime verification. |
| [Git](https://git-scm.com/downloads) | current | Source control and generated-output checks. |
| [`vx`](https://github.com/vandordev/vx) | latest | Internal scaffold runtime called only by `task add:*`. |

On macOS or Linux, install Bun, Task, RTK, and `vx` with:

```bash
curl -fsSL https://bun.com/install | bash
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b "$HOME/.local/bin"
curl -fsSL https://raw.githubusercontent.com/rtk-ai/rtk/refs/heads/master/install.sh | sh
go install github.com/vandordev/vx/cmd/vx@latest
```

Install Go, Docker, and Git using their official installer or your operating
system package manager. Add `$HOME/.local/bin`, `$HOME/go/bin`, and Bun's
installation directory to `PATH` if the installer did not do so. Air requires
no global installation: `task install` downloads it from the `go.mod` tool
directive.

Verify the setup and bootstrap a local checkout:

```bash
task install
cp .env.example .env
# Set development-safe values for the required secrets in .env.
task doctor
docker compose up -d db
task migrate
```

`task doctor` reports missing required tools or local files. It also verifies
that the checked-in Go configuration can load. Do not call `vx` or Air
directly in normal project workflows; Taskfile tasks own both integrations.

## Commands

```bash
task install
task doctor
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
task test:go
task test:web
task test:realtime
task test:config
task quality
task build
task ci
```

`task dev` runs API, worker, scheduler, and web together. API, worker, and
scheduler use the project-scoped Air tool for Go hot reload; `apps/migrate`
remains a one-shot command. Select services with arguments after `--`, for
example `task dev -- api web realtime`; realtime is not part of the default
set. Run `task migrate` separately before starting the long-lived services.

Run `task doctor` after setup to verify local tools, Docker Compose, `.env`,
and Go configuration loading. Use the segmented `task test:*` commands for
focused feedback. `task ci` runs `task quality` followed by `task build`, the
same full verification expected before integration.

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

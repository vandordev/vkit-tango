# Go Backend, River, and Shared Configuration Design

**Status:** Proposed — approved architectural direction, awaiting implementation-plan review  
**Date:** 2026-07-18

## 1. Overview

The repository becomes a mixed Go and TypeScript monorepo with a single Go
write-side backend. Go owns HTTP mutations, background jobs, scheduling, and
all PostgreSQL mutations. TypeScript remains responsible for the TanStack Start
web application and the Socket.IO realtime runtime.

The central rule is deliberately simple: every domain mutation goes through a
Go usecase. This prevents synchronous API code and asynchronous worker code
from developing separate business rules or separate write paths.

PostgreSQL is the system database and River is the durable queue. Redis,
Asynq, Prisma, Elysia, pg-boss, the standalone TypeScript worker, and the
standalone TypeScript scheduler are not part of the target architecture.

This design supersedes earlier architecture documents where they prescribe
Prisma, Elysia, pg-boss, a TypeScript worker/scheduler, or uppercase
environment-first configuration. Historical documents remain intact as records
of the previous baseline.

The committed module targets Go 1.25. Ent's pinned code-generation dependency
set requires this minimum; the repository's installed Go toolchain is 1.25.7.

## 2. Goals and Non-Goals

### Goals

- Establish one source of truth for every business mutation in Go.
- Keep reads lightweight: HTTP handlers and job handlers may query Ent directly
  when they only shape a read model.
- Use PostgreSQL for both application data and durable River jobs, avoiding a
  Redis dependency.
- Make HTTP contracts versioned from day one under `/api/v1/*`.
- Generate OpenAPI from the Go HTTP API and generate TypeScript SDK and
  TanStack Query support from that specification.
- Keep Socket.IO realtime while ensuring the browser treats events as a signal
  to refetch canonical HTTP state.
- Use one snake_case YAML configuration source for Go and TypeScript runtimes
  without coupling their startup processes.

### Non-Goals

- Introduce full DDD layers, repositories per table, operation classes,
  duplicated entities, or an Fx-style dependency registry.
- Make River job payloads a cross-language contract. API and workers are Go;
  only realtime events cross the Go/TypeScript boundary.
- Put product-specific models, authentication rules, or example domains in the
  reusable baseline.
- Provide exactly-once external side effects. River gives durable,
  at-least-once processing; integrations must remain idempotent.

## 3. System Context

```text
Browser
  │ HTTP: same-origin /api/v1/*                         Socket.IO /ws
  ▼                                                        ▼
TanStack Start (TypeScript)                         Realtime (TypeScript)
  │ proxy in development                                  ▲
  ▼                                                        │ authenticated internal event
Huma API (Go) ─────── Ent ───────► PostgreSQL ◄──── River worker + scheduler (Go)
  │                         ▲          │                         │
  └─ inserts River job in same tx ─────┘          uses the same Go usecases
                                                        │
                                                        └── external integrations
```

In production, the edge proxy sends `/` to `apps/web`, `/api/*` and `/health`
to `apps/api`, and the Socket.IO path (normally `/ws`) to `apps/realtime`.
During development, Vite proxies `/api/*` and `/health` to the Go API. TanStack
Start must not embed or duplicate the Go API implementation.

## 4. Runtime and Repository Boundaries

The root is both a Go module and a JavaScript workspace. Mixing the languages
is intentional: a Go application under `apps/` can import root `internal/`
packages, while TypeScript workspaces retain their usual boundaries.

```text
apps/
  web/                 # TanStack Start; presentation and generated API client
  api/                 # Go Huma process entrypoint
  worker/              # Go River worker process entrypoint and periodic schedules
  realtime/            # TypeScript Socket.IO process entrypoint
  migrate/             # Go Goose migration process entrypoint

database/
  schema/              # Ent schema source of truth for the current data model
  migrations/          # immutable Goose migration history

internal/
  bootstrap/           # explicit dependency composition for Go runtimes
  config/              # Go YAML loader, typed config, validation
  usecase/             # all intent-named domain mutations
  transport/http/      # Huma route registration, DTOs, version bundles
  worker/river/        # River job definitions, handlers, periodic registration
  platform/
    db/                # generated Ent client; never hand-edited
    river/             # River client and queue integration
    realtime/          # authenticated HTTP publisher for Socket.IO runtime
    external/          # integration clients behind narrow interfaces
    observability/     # logging, metrics, tracing setup

config/                # shared snake_case YAML configuration and test fixtures
contracts/
  openapi/             # generated, reviewed OpenAPI document
  asyncapi/            # versioned realtime event contract
```

`internal/context` is intentionally absent. `internal/usecase` is a pragmatic
write-side boundary, not a mandate to model the project as a large DDD system.
`internal/platform/db` is generated from `database/schema`; code outside that
directory consumes the generated Ent client but does not modify it.

## 5. Persistence, Migrations, and Ent

Ent is the sole application ORM. It is used for both mutations and direct
reads. No Go code uses Prisma or sqlc.

- `database/schema` defines the current Ent schema.
- `SystemMetadata` is the sole platform-owned baseline entity. It stores
  non-product metadata by unique key and uses an application-generated UUID
  primary key; it must not become a generic product-data table.
- `database/migrations` stores ordered Goose migrations and is the only means
  of changing a deployed database.
- `apps/migrate` executes Goose. Production never calls `ent.Schema.Create`.
- Ent generation targets `internal/platform/db` and is invoked by a reproducible
  task; generated output is committed or otherwise handled consistently by the
  repository's generation policy, but is never edited manually.
- River's PostgreSQL migrations are incorporated in the same Goose migration
  flow and pinned to the selected River version. River tables such as
  `river_job`, `river_leader`, and `river_migration` are therefore created and
  upgraded with the application schema.

Every Go process receives the same explicitly composed Ent client. Transaction
ownership belongs to the caller that changes business state, normally a
usecase. River insertion must occur in that same transaction when a mutation
requires background work.

## 6. Mutation and Read Discipline

### Mutations

Only `internal/usecase` may invoke Ent write operations (`Create`, `Update`,
`Delete`, upsert, or equivalent) against product data. A usecase:

1. accepts an intent-specific input, such as `CancelBooking` rather than a
   generic status update;
2. enforces business rules and idempotency;
3. owns the Ent transaction;
4. changes domain state;
5. inserts any required River jobs in the same transaction; and
6. returns a result suitable for its caller to map to a transport DTO.

HTTP mutation handlers validate transport input and invoke one usecase. River
handlers decode job arguments, perform any needed reads, and invoke the same
usecase for a write. This avoids duplicate business rules between synchronous
and asynchronous paths.

### Reads

Read paths may use Ent directly from a Huma handler, River handler, or another
caller that needs a read model. They may select only the needed fields and map
them to a transport-specific response. Ent entities are never exposed as the
HTTP contract.

Narrow interfaces are appropriate at external boundaries (payments, mail,
storage, external HTTP, and realtime). Repository interfaces or wrappers around
every Ent table are not part of this baseline.

## 7. HTTP API and Versioning

Huma is the only public HTTP API server. The public same-origin API prefix
remains `/api`, and every active business endpoint is versioned:

```text
/api/v1/<resource-or-action>
```

`/health` remains an unversioned process-health endpoint. API documentation and
its OpenAPI artifact are exposed at stable non-business locations such as
`/api/docs` and `/api/openapi.json`.

`internal/transport/http/method` provides the only route-registration wrapper
available to feature packages. It binds an explicit API version and constructs
the full path. Feature routes supply a relative path only. The wrapper rejects
paths that already contain `/api` or a version segment, so an accidental
unversioned or double-prefixed handler cannot be registered.

```text
internal/transport/http/
  method/
  v1/
    routes.go
    <feature>/
  v2/                 # added only for a breaking public contract
```

Operations use explicit, stable, versioned IDs such as `v1_create_booking`.
Breaking behavior receives a new version bundle; v1 is not silently changed.

## 8. River Jobs and Scheduling

River replaces Redis/Asynq and pg-boss. It uses the existing PostgreSQL
database for durable job state.

- The Go API inserts jobs with River's transactional insertion API inside the
  owning Ent transaction.
- The Go worker processes jobs. Where processing mutates the database, it uses
  the shared usecase and completes the job transactionally with that mutation.
- Job payloads and side effects are idempotent. Use unique-job options and
  domain idempotency keys where appropriate; an at-least-once job may run more
  than once.
- The worker also registers River periodic jobs. There is no `apps/scheduler`
  process.
- All worker replicas register the same periodic schedule. River leader
  election ensures one scheduler is active, but periodic jobs are not treated
  as a precise, durable clock.
- Any business action that must not be missed uses a reconciliation job that
  scans for due database rows and applies an idempotent usecase. A periodic tick
  only triggers that reconciliation.

The first implementation must prove, with a PostgreSQL integration test, that
the chosen Ent driver and River client can share a transaction for both job
insertion and transactional job completion. This is a hard acceptance check,
not an assumption.

## 9. Realtime and AsyncAPI

`apps/realtime` remains the TypeScript Socket.IO service. Go processes do not
implement the Socket.IO protocol. Instead, `internal/platform/realtime` sends
an authenticated internal HTTP request to a private realtime endpoint (for
example `/internal/events`). The endpoint validates the internal credential and
emits the corresponding Socket.IO event to the intended room(s).

When a usecase commits state that needs a browser update, it inserts a
`realtime.publish.v1` River job in the same database transaction. The job
handler calls the internal realtime endpoint. River retries publication; a
duplicate event is safe because the browser uses it only as an invalidation
signal and refetches canonical state from `/api/v1/*`.

AsyncAPI documents this Go-to-TypeScript realtime event boundary, not River job
arguments. Events are versioned and carry enough routing and invalidation data
for the web client, while the API remains the source of truth. The contract is
the basis for compatible Go and TypeScript event types or validation code.

## 10. Configuration

`config/` is the shared, checked-in configuration source for both languages.
Files use snake_case keys and semantic roots:

```text
config/
  app.yaml
  database.yaml
  http_api.yaml
  worker.yaml
  realtime.yaml
  web.yaml
  observability.yaml
  testdata/
  README.md
```

For example, `http_api.yaml` begins with `http_api:`, and `web.yaml` contains a
`web.public` subtree. Go's `internal/config` and TypeScript's `packages/config`
load the same YAML files independently, then validate their own typed subset.
There is no configuration compiler and no runtime dependency between the two
loaders.

The supported scalar interpolation syntax is deliberately limited to
`${NAME}` and `${NAME:-fallback}`. Interpolation happens before typed decoding;
recursive expansion and implicit environment-variable reads are not supported.
Secrets are never committed as YAML literals. Local `.env` files and deployment
secrets provide values only for placeholders explicitly present in YAML.

Runtime ownership is as follows:

| Runtime | YAML modules |
| --- | --- |
| Go API | app, database, http_api, realtime, observability |
| Go worker | app, database, worker, realtime, observability |
| Go migrate | database |
| TypeScript realtime | app, realtime, observability |
| TypeScript web | app, web |

The browser receives only `web.public`. A small TypeScript adapter may expose
that safe subset through Vite `VITE_*` values when required, but `VITE_*` is not
written into shared YAML and server secrets can never enter browser config.

Both loaders use the same fixtures in `config/testdata` to verify required
variables, defaults, invalid values, and public-config exclusion.

## 11. OpenAPI and TypeScript Client Generation

Huma generates the authoritative OpenAPI 3.1 document from registered Go
routes and DTOs. The generated artifact is stored at
`contracts/openapi/openapi.json` for review and code generation.

Hey API generates the web client from that local artifact:

```text
contracts/openapi/openapi.json
  └── Hey API
        └── apps/web/src/lib/api/generated/
              types/
              sdk/
              tanstack-query/
```

The initial generator uses the `@hey-api/client-fetch`, `@hey-api/sdk`, and
`@tanstack/react-query` plugins, with exact package versions pinned. Generated
files are never edited. A hand-written `apps/web/src/lib/api/client.ts` configures
the shared same-origin client, including request-scoped concerns such as
credentials or headers. Server-side web code uses the generated SDK; browser
code uses generated TanStack Query options/hooks.

Taskfile commands provide a reproducible workflow:

| Task | Responsibility |
| --- | --- |
| `task api:openapi` | Build/export the Huma OpenAPI document without starting a listener. |
| `task api:client:generate` | Export OpenAPI then run Hey API generation. |
| `task api:client:check` | Regenerate and fail if the working tree differs. |

Build and CI run the stale-generated-client check whenever API contracts or web
client generation are in scope.

## 12. Migration Scope and Sequencing

Implementation proceeds as a coherent platform migration rather than a partial
dual-write deployment:

1. Add root Go module, target directories, shared configuration fixtures, Ent
   schema/code generation, Goose runner, and River migrations.
2. Establish explicit Go bootstrap, Huma API/version wrapper, OpenAPI export,
   and Hey API generation before migrating feature endpoints.
3. Move each existing mutation to a Go usecase and Huma v1 handler; replace
   TypeScript consumers with the generated client.
4. Move each job handler and periodic schedule to Go River. Remove the
   TypeScript worker and scheduler only after their workloads have moved.
5. Connect transactional realtime publication and replace obsolete queue/API
   configuration and dependencies.
6. Remove Elysia, Prisma, pg-boss, Redis/Asynq remnants, old process entry
   points, unused Compose services, and their environment schemas once no
   runtime references them.

The migration does not leave two mutation implementations for the same domain
operation. A feature switches only when its API mutation and its worker path
both invoke the same Go usecase.

## 13. Operational and Security Requirements

- Every process has graceful shutdown: stop accepting HTTP work, drain or stop
  River according to its documented shutdown behavior, then close database
  resources.
- Use structured logs with request IDs, job IDs, attempt count, API version,
  and domain identifiers safe for logs.
- Apply timeouts to database calls, internal realtime calls, and external
  integrations. Configure retry policy by job type; do not retry permanent
  validation failures.
- The internal realtime endpoint is not publicly exposed through the edge
  proxy. It authenticates Go publishers with a dedicated secret held only by
  trusted server runtimes.
- All browser state-changing requests use the Go API and its authentication/CSRF
  policy; the precise product authentication model remains outside this generic
  architecture migration.

## 14. Verification and Acceptance Criteria

The implementation is acceptable when all of the following are demonstrable:

- A v1 HTTP mutation calls a Go usecase, which writes through Ent; no TypeScript
  code mutates product PostgreSQL data.
- The same usecase is used by an HTTP handler and a River job handler.
- A mutation and its River enqueue commit or roll back together in an actual
  PostgreSQL integration test.
- A job mutation and transactional completion commit or roll back together in
  an actual PostgreSQL integration test.
- A failed realtime delivery is retried by River; duplicate browser events only
  cause harmless refetches.
- Periodic jobs are registered by `apps/worker`, and important time-based work
  is protected by a reconciliation query rather than a one-shot tick.
- Fresh PostgreSQL setup succeeds using Goose only, including River tables.
- All active business endpoints are under `/api/v1/*`; the registration wrapper
  rejects an unversioned feature path.
- `task api:openapi`, `task api:client:generate`, and
  `task api:client:check` pass, and the web uses generated SDK/query artifacts.
- Go and TypeScript configuration tests pass against the shared YAML fixtures,
  and browser configuration cannot contain database or realtime internal
  secrets.

## 15. Risks and Explicit Follow-ups

| Item | Treatment |
| --- | --- |
| Ent and River transaction compatibility | Prove it in the first PostgreSQL integration test before migrating business workloads. |
| River periodic leader handover can miss a tick | Use database-backed reconciliation for important deadlines. |
| At-least-once jobs can duplicate external work | Require idempotency keys and record durable domain state before external calls where applicable. |
| Node Socket.IO ticket/session compatibility | Preserve the current runtime contract during migration; specify the Go-issued ticket/auth integration with the product auth design before enabling protected rooms. |
| Huma/OpenAPI and Hey API generator changes | Pin versions, commit generated contracts, and enforce regeneration in CI. |
| Broad migration scope | Migrate feature-by-feature but never maintain two mutation implementations for one operation. |

## 16. References

- [River TypeScript documentation](https://riverqueue.com/docs/typescript)
- [River transactional job completion](https://riverqueue.com/docs/transactional-job-completion)
- [River periodic jobs](https://riverqueue.com/docs/periodic-jobs)
- [River migrations](https://riverqueue.com/docs/migrations)
- [Ent code generation](https://entgo.io/docs/code-gen/)
- [Ent transactions](https://entgo.io/docs/transactions/)
- [Goose](https://github.com/pressly/goose)
- [Huma](https://github.com/danielgtaylor/huma)
- [Hey API](https://heyapi.dev/)

# Fx Command and Generator Design

**Status:** Approved design — awaiting implementation-plan review  
**Date:** 2026-07-19

## 1. Purpose

Introduce explicit dependency injection with `go.uber.org/fx`, an
intent-oriented command pattern, and a small `vpkg`/`vx` generator surface for
the Go backend. The design removes manual runtime wiring from application
entrypoints while retaining the baseline's deliberately small scope.

This design supersedes earlier architecture statements that reject an Fx
registry or an `apps/scheduler` process. It does not add product-specific
domains, repositories, or DDD layers.

## 2. Goals and Non-Goals

### Goals

- Make `apps/api`, `apps/worker`, and `apps/scheduler` thin Fx composition
  roots with explicit runtime-specific modules.
- Put every business mutation behind a reusable command in `internal/usecase`.
- Ensure HTTP handlers, River jobs, and schedules call or enqueue commands
  without copying business logic.
- Keep command tests as the authoritative behavioral test suite.
- Provide one developer-facing Taskfile command for every scaffold and sync
  action; developers and CI do not invoke `vx` directly.
- Generate deterministic Fx and registration registries from Go source.

### Non-Goals

- Generate a handler, job, or schedule automatically for every use case.
- Add a test file for every HTTP handler, job, or schedule adapter.
- Replace Ent, Goose, River, Huma, or the existing shared YAML configuration.
- Use reflection at application runtime to discover dependencies or routes.
- Turn `apps/migrate` into an Fx runtime; it remains a short-lived migration
  command.

## 3. Runtime Composition

```text
apps/api       = shared module + commands + HTTP module
apps/worker    = shared module + commands + River job module
apps/scheduler = shared module + scheduler module + River enqueue client
apps/migrate   = Goose/River migration process; no Fx
```

Each long-lived app starts with `fx.New(Module).Run()`. The root's only
responsibility is selecting modules for that process; it contains no manual
construction of a database, producer, use case, handler, worker, or schedule.

Shared infrastructure modules provide typed configuration, PostgreSQL/Ent,
River producer/client primitives, logging, observability, and external
clients. Every resource with a connection or server lifecycle is started and
stopped through `fx.Lifecycle`.

The scheduler is a separate process because it owns only periodic enqueueing
and reconciliation triggering. It never performs domain mutations directly.
The worker owns execution of River jobs. Both must be safe when horizontally
replicated and rely on River's supported periodic scheduling/leader behavior.

## 4. Shared Contracts and Commands

`internal/contract` holds reusable interfaces shared by commands and their
adapters. Feature packages must not repeat one-off interfaces merely to expose
a command to HTTP, a job, or a scheduler.

```text
internal/contract/
  command.go     # Command[I, O] with Execute(context.Context, I) (O, error)
  http.go        # HTTPHandler with RegisterRoutes()
  job.go         # River worker-registration boundary
  scheduler.go   # periodic-registration boundary
```

Every use case is one public struct in one file and implements an exact
`contract.Command[Input, Output]`. Its constructor returns the concrete
pointer, but generated Fx registration exposes it as the matching generic
contract using `fx.Annotate(..., fx.As(...))`. Therefore an adapter depends on
the command contract rather than a concrete use-case implementation.

```text
internal/usecase/
  set_system_metadata.go
  set_system_metadata_test.go
```

The use-case test is mandatory and contains business rules, transaction
behavior, idempotency, error cases, and side effects. It is the only mandatory
test companion imposed by the file convention:

```text
1 use-case file = 1 public use-case struct = 1 use-case test file
```

HTTP handlers, River jobs, and schedules are deliberately thin adapters. They
have one public adapter struct per operation/registration file but no required
per-file unit test. Transport contract or integration tests may cover routing,
serialization, error mapping, and OpenAPI when those behaviors warrant it.

## 5. HTTP: Chi, Huma, and One Operation per Handler

HTTP uses Chi as the router and `huma/v2/adapters/humachi` as the Huma adapter.
The HTTP module creates a Chi router, installs cross-cutting middleware, builds
the Huma API with `humachi.New`, and exposes the API as a versioned group.
Health, readiness, OpenAPI, and docs remain process-level routes according to
the repository's API conventions.

`internal/transport/http/method` is the sole Huma registration wrapper used by
feature handlers. It supplies typed helpers for `GET`, `POST`, `PUT`, `PATCH`,
and `DELETE`, sets stable operation metadata, and delegates to
`huma.Register`.

One HTTP operation means one HTTP method plus one path. Different methods on
the same path are separate handlers. For example, `GET /api/v1/me` and
`PATCH /api/v1/me` have distinct files and distinct handler structs.

```text
internal/transport/http/handler/system_metadata/
  set_system_metadata.go  # SetSystemMetadataHandler; method.PUT(...)
```

A handler validates and maps HTTP input, calls one `contract.Command`, and
maps the result or domain error to an HTTP response. It does not contain
transactions, persistence mutations, River insertion, or business policy.

## 6. River Jobs and Schedules

`internal/worker/river` contains typed River jobs. A job decodes its typed
arguments and calls the relevant `contract.Command`; it does not restate the
command's mutation logic. When suitable, a job argument may reuse the
command's serializable input type to avoid an artificial mapping layer.

`internal/scheduler/river` contains periodic registrations. A schedule only
enqueues a typed job or reconciliation trigger. It must not invoke a use case
directly. Deadline-sensitive business work uses idempotent reconciliation
jobs; a periodic tick is never the source of truth.

Generated worker and scheduler registries install only types that meet the
shared `internal/contract` registration interfaces. The API root excludes job
and periodic registration; the worker root excludes HTTP registration; the
scheduler root excludes HTTP and job execution.

## 7. Generator Layout and Responsibilities

The repository gains a small, versioned package at `vpkg/vandor/go`. It uses
the `vx` preview-first template runtime for scaffolding and small Go tools for
source-aware synchronization.

```text
vpkg/vandor/go/
  vpkg.yaml
  templates/
    usecase.vxt
    http_handler.vxt
    job.vxt
    scheduler.vxt
  tools/
    sync/
```

`vx` is responsible only for creating new files. The generated templates use
the injected Go project context, derive identifiers and paths from one `name`
input, and refuse to overwrite existing files. The use-case template creates
the required `*_test.go`; adapter templates do not.

The Go sync tool parses non-test Go source with the Go AST, validates required
constructors and shared-contract conformance, sorts all output deterministically,
formats it with `go/format`, and writes generated code to a clearly marked
directory such as `internal/generated/fx`. Generated files are never edited
by hand.

The sync tool generates four registries:

```text
internal/generated/fx/
  usecases_gen.go   # fx.Provide + fx.As(command contract)
  http_gen.go       # handler providers and route registration invokes
  worker_gen.go     # typed River worker registration
  scheduler_gen.go  # periodic registration
```

Runtime discovery through reflection is prohibited. The scanner recognizes
only documented directories, exported `New...` constructors, and the shared
contract boundaries, so a missing or invalid adapter fails at sync/build time.

## 8. Taskfile Developer Interface

Taskfile is the only public developer and CI interface. Individual add tasks
call `vx gen` internally and then run only the synchronization needed by the
new surface. `task sync` regenerates and verifies every registry.

```sh
task add:usecase name=SetSystemMetadata

task add:http-handler \
  name=SetSystemMetadata \
  method=PUT \
  path=/api/v1/system-metadata/{key}

task add:job name=SetSystemMetadata
task add:scheduler name=ReconcileSystemMetadata

task sync:usecase
task sync:http
task sync:worker
task sync:scheduler
task sync
```

The `name` in an adapter task identifies the related use case. An HTTP handler
does not need a duplicate `command=` argument: the template imports the
corresponding input/output types and receives the matching generic command
contract. The HTTP `method` value determines the `method.GET`, `method.POST`,
`method.PUT`, `method.PATCH`, or `method.DELETE` registration call. The full
versioned `path` is explicit in the task so there is no hidden route inference.

`task add:job` is the canonical job command; `river-job` is not exposed in the
Taskfile API. River remains the implementation detail of the job adapter.

## 9. Error Handling and Verification

- Task validation rejects missing or unsupported scaffold values before files
  are created.
- `vx` previews by default; Taskfile's add commands pass `--apply` only after
  all required values are present. Existing planned paths fail rather than
  being overwritten.
- A failed targeted sync leaves no hand-edited registry; re-running the same
  sync is safe and deterministic.
- `task sync` must be followed by `gofmt` and a focused Go test/build check.
- Shared or runtime changes additionally run `task quality` and `task build`.
- Generated output is inspected in version control; CI fails when a sync would
  change committed generated files.

## 10. Acceptance Criteria

1. `go.mod` includes Fx and all long-lived Go runtimes use Fx modules.
2. `apps/api`, `apps/worker`, and `apps/scheduler` have distinct composition
   roots; `apps/migrate` remains short-lived and non-Fx.
3. A generated use-case registry exposes every valid use case through
   `contract.Command[I, O]` without hand-maintained provider lists.
4. Chi + `humachi` serves Huma routes, and all generated handler scaffolds use
   the shared typed method helpers.
5. A `PUT` handler can be added with only `name`, `method`, and `path` and is
   registered after the targeted sync.
6. Jobs invoke commands and schedules enqueue jobs; neither repeats a
   mutation's business rules.
7. Every scaffold and sync operation is available through Taskfile commands;
   no documented workflow requires calling `vx` directly.
8. Every generated use case includes its matching test file; generated
   adapters do not require per-file tests.

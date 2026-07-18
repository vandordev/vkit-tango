# Go Backend, River, and Shared Configuration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the TypeScript write-side stack with a Go Huma/Ent/River backend while retaining TanStack Start and Socket.IO, with versioned HTTP, YAML configuration, and generated TypeScript API clients.

**Architecture:** The root becomes a Go module alongside the Bun workspace. Go API, worker, migration, usecase, Ent, River, and configuration code live below `apps/`, `internal/`, and `database/`; TypeScript owns web presentation and Socket.IO. Ent mutations and River enqueue/completion share PostgreSQL transactions, while realtime delivery is a retryable River job to a private Socket.IO HTTP endpoint.

**Tech Stack:** Go 1.24, Huma v2, Ent, pgx v5, River, Goose, PostgreSQL 16, TanStack Start, Socket.IO, AsyncAPI, Hey API, Bun, Taskfile, Docker Compose.

---

## Preconditions and implementation rules

- Read [the approved architecture spec](../specs/2026-07-18-go-backend-river-architecture-design.md) before every phase; it is the source of truth if an older repository file disagrees.
- Work in an isolated worktree when executing this plan. Do not start the migration in a dirty shared tree.
- Do not introduce a sample product entity, product usecase, or product job. The baseline only provides the platform and test harnesses that a product feature will use.
- Do not leave Elysia, Prisma, pg-boss, Eden, a TypeScript scheduler, or an embedded API path reachable after the final cleanup task.
- Each task below ends with its stated focused verification before proceeding. Commit the files listed in that task with the supplied message.

## Target file map

| Path | Responsibility |
| --- | --- |
| `go.mod`, `go.sum` | Root Go module and pinned backend dependencies. |
| `database/schema/entc.go`, `database/schema/schema.go` | Ent generation configuration and empty domain-neutral schema package. |
| `database/migrations/*.sql` | Immutable Goose migrations, including River vendor migrations. |
| `internal/platform/db/` | Generated Ent client and database connection helpers. |
| `internal/config/` | Go YAML loading, interpolation, typed validation, and shared-fixture tests. |
| `internal/bootstrap/` | Explicit construction of DB, River, realtime publisher, API, and worker dependencies. |
| `internal/transport/http/` | Huma API construction, response envelope, version-safe routing, v1 routes, and OpenAPI export. |
| `internal/usecase/` | Intent-named mutation boundary and transaction helper; initially no product usecase. |
| `internal/worker/river/` | Typed River jobs, transactional job completion, periodic reconciliation registration, and tests. |
| `internal/platform/realtime/` | Go authenticated HTTP publisher for the Socket.IO runtime. |
| `apps/api/main.go` | Go API process entrypoint. |
| `apps/worker/main.go` | Go River worker entrypoint and periodic registration. |
| `apps/migrate/main.go` | Goose migration entrypoint. |
| `apps/realtime/` and `packages/realtime/` | Socket.IO runtime and generated/validated realtime event contract. |
| `apps/web/src/lib/api/` | Hey API configuration, generated outputs, typed client wrapper, and query integration. |
| `apps/web/vite.config.ts`, `apps/web/src/server.ts` | Development API proxy and removal of the Elysia embedding adapter. |
| `Taskfile.yml`, `package.json`, `docker-compose.yml`, `Dockerfile.*` | Unified developer, generation, migration, and container commands. |
| `README.md`, `.agent/*.md`, `config/README.md`, `.env.example` | Updated public and contributor contract. |

### Task 1: Establish the Go module and executable skeleton

**Files:**
- Create: `go.mod`
- Create: `go.sum`
- Create: `apps/api/main.go`
- Create: `apps/worker/main.go`
- Create: `apps/migrate/main.go`
- Create: `internal/bootstrap/bootstrap.go`
- Create: `internal/bootstrap/bootstrap_test.go`
- Modify: `.gitignore`
- Test: `internal/bootstrap/bootstrap_test.go`

- [ ] **Step 1: Write the failing composition test.**

  Define a bootstrap API that accepts already-constructed dependencies, so tests
  do not need PostgreSQL or environment variables:

  ```go
  package bootstrap_test

  func TestNewRejectsNilDatabase(t *testing.T) {
      _, err := bootstrap.New(bootstrap.Dependencies{})
      require.ErrorIs(t, err, bootstrap.ErrDatabaseRequired)
  }
  ```

- [ ] **Step 2: Run the focused test and record the initial failure.**

  Run: `rtk go test ./internal/bootstrap -run TestNewRejectsNilDatabase -count=1`  
  Expected: FAIL because the Go module and `bootstrap` package do not exist.

- [ ] **Step 3: Create the root Go module and pinned dependencies.**

  Initialize the module as `github.com/vandordev/vkit-fast` and add exact
  compatible versions of these direct dependencies: `entgo.io/ent`,
  `github.com/danielgtaylor/huma/v2`, `github.com/jackc/pgx/v5`,
  `github.com/pressly/goose/v3`, `github.com/riverqueue/river`,
  `github.com/stretchr/testify`, and `gopkg.in/yaml.v3`. Keep the Go directive
  at `go 1.24` and let `go mod tidy` own `go.sum`.

- [ ] **Step 4: Implement explicit bootstrap validation and process skeletons.**

  `internal/bootstrap/bootstrap.go` exposes only typed dependencies and does
  not use a service locator:

  ```go
  package bootstrap

  var ErrDatabaseRequired = errors.New("database is required")

  type Dependencies struct { Database *ent.Client }
  type Runtime struct { Database *ent.Client }

  func New(deps Dependencies) (*Runtime, error) {
      if deps.Database == nil { return nil, ErrDatabaseRequired }
      return &Runtime{Database: deps.Database}, nil
  }
  ```

  Add `main.go` programs with `signal.NotifyContext` for `SIGINT` and `SIGTERM`.
  They may only construct their runtime and wait for cancellation at this stage;
  they must not start a dummy server, mutate schema, or contain product logic.
  Add `.air/`, `tmp/`, and Go coverage output to `.gitignore` without ignoring
  generated Ent code.

- [ ] **Step 5: Run formatting, focused tests, and module verification.**

  Run:

  ```bash
  rtk gofmt -w apps internal
  rtk go mod tidy
  rtk go test ./internal/bootstrap -count=1
  rtk go vet ./...
  ```

  Expected: all commands exit 0.

- [ ] **Step 6: Commit the skeleton.**

  ```bash
  git add go.mod go.sum .gitignore apps/api/main.go apps/worker/main.go apps/migrate/main.go internal/bootstrap
  git commit -m "build: add go backend runtime skeleton"
  ```

### Task 2: Replace the configuration model with shared snake_case YAML

**Files:**
- Create: `config/app.yaml`
- Create: `config/database.yaml`
- Create: `config/http_api.yaml`
- Create: `config/observability.yaml`
- Create: `config/testdata/required.yaml`
- Create: `config/testdata/defaults.yaml`
- Create: `config/testdata/invalid.yaml`
- Create: `config/README.md`
- Create: `internal/config/loader.go`
- Create: `internal/config/config.go`
- Create: `internal/config/loader_test.go`
- Create: `internal/config/config_test.go`
- Modify: `config/web.yaml`
- Modify: `config/worker.yaml`
- Modify: `config/realtime.yaml`
- Modify: `packages/config/src/loader.ts`
- Modify: `packages/config/src/index.ts`
- Modify: `packages/config/src/run.ts`
- Modify: `packages/config/src/loader.test.ts`
- Modify: `packages/config/src/run.test.ts`
- Modify: `packages/config/src/realtime.ts`
- Modify: `packages/config/src/realtime.test.ts`
- Delete: `packages/config/src/api.ts`
- Delete: `packages/config/src/api.test.ts`
- Delete: `packages/config/src/common.ts`
- Delete: `packages/config/src/common.test.ts`
- Delete: `packages/config/src/deployment.test.ts`
- Delete: `packages/config/src/scheduler.ts`
- Delete: `packages/config/src/scheduler.test.ts`
- Delete: `packages/config/src/storage.ts`
- Delete: `packages/config/src/storage.test.ts`
- Modify: `.env.example`
- Delete: `config/base.yaml`
- Delete: `config/api.yaml`
- Delete: `config/scheduler.yaml`
- Delete: `config/storage.yaml`
- Test: `internal/config/loader_test.go`, `internal/config/config_test.go`, `packages/config/src/loader.test.ts`, `packages/config/src/run.test.ts`, `packages/config/src/realtime.test.ts`

- [ ] **Step 1: Write Go and TypeScript failing tests against the same fixtures.**

  Both test suites must resolve these exact scalar forms:

  ```yaml
  database:
    url: ${DATABASE_URL}
  http_api:
    port: ${API_PORT:-4101}
  ```

  Assert all of the following in each language: a present `DATABASE_URL` is
  substituted; absent `API_PORT` becomes `4101`; absent required
  `DATABASE_URL` returns an error naming that variable; `$DATABASE_URL` is a
  literal string; nested `${A:-${B}}` is rejected; and a web runtime cannot
  obtain `database.url` or `realtime.internal_api_key`.

- [ ] **Step 2: Run both focused suites to prove the old configuration contract fails.**

  Run:

  ```bash
  rtk go test ./internal/config -count=1
  rtk bun test packages/config
  ```

  Expected: Go fails because the package is absent; Bun fails after the fixture
  assertions are added because the old uppercase flat module format differs.

- [ ] **Step 3: Define the final YAML modules and public boundary.**

  Use these semantic roots and no uppercase configuration keys:

  ```yaml
  # config/http_api.yaml
  http_api:
    host: 0.0.0.0
    port: ${API_PORT:-4101}
    public_base_url: ${API_PUBLIC_BASE_URL:-http://localhost:4101}

  # config/web.yaml
  web:
    public:
      api_base_url: ${WEB_API_BASE_URL:-http://localhost:4101}
      realtime_url: ${WEB_REALTIME_URL:-http://localhost:4102}
  ```

  Put database URL under `database.url`; place the Socket.IO publishing secret
  only under `realtime.internal_api_key`; place the ticket secret only under
  `realtime.ticket_secret`. Keep all secret values interpolation-only.

- [ ] **Step 4: Implement the two independent typed loaders.**

  Go exposes a single explicit entrypoint:

  ```go
  type Loader struct { Directory string; Environment map[string]string }
  func (l Loader) Load(modules ...string) (map[string]any, error)
  func LoadAPI(l Loader) (API, error)
  func LoadWorker(l Loader) (Worker, error)
  func LoadMigrate(l Loader) (Migrate, error)
  ```

  The loader parses YAML, walks only scalar strings for `${NAME}` and
  `${NAME:-fallback}`, rejects any string containing an unrecognized `${` form,
  then decodes and validates the selected typed struct. TypeScript keeps the
  same `loadConfig` name but accepts only the new module names and returns
  nested snake_case records. Replace environment-schema factories with typed
  selectors for `web` and `realtime`; no feature code reads `process.env`.

- [ ] **Step 5: Add the browser adapter and eliminate old modules.**

  Export `createWebPublicConfig` from `packages/config`, returning only
  `web.public`. Make `packages/config/src/run.ts` set only the explicit Vite
  values derived from this object. Delete old `base`, `api`, `scheduler`, and
  `storage` modules and tests tied to uppercase `PORT`, `DATABASE_URL`, or
  module deep-merge behavior. Update `.env.example` to list only interpolation
  variable names used by the new YAML modules.

- [ ] **Step 6: Run configuration verification.**

  Run:

  ```bash
  rtk gofmt -w internal/config
  rtk go test ./internal/config -count=1
  rtk bun test packages/config
  rtk bunx prettier --check "config/**/*.yaml" "config/**/*.md"
  ```

  Expected: each suite passes its required/default/invalid/public-exclusion
  cases and formatting exits 0.

- [ ] **Step 7: Commit the shared configuration foundation.**

  ```bash
  git add config internal/config packages/config .env.example
  git rm config/base.yaml config/api.yaml config/scheduler.yaml config/storage.yaml
  git commit -m "feat: add shared yaml runtime configuration"
  ```

### Task 3: Add Ent, Goose, and River migrations with transaction proof

**Files:**
- Create: `database/schema/entc.go`
- Create: `database/schema/schema.go`
- Create: `database/migrations/00001_river.sql`
- Create: `database/migrations/00002_ent_schema.sql`
- Create: `internal/platform/db/open.go`
- Create: `internal/platform/db/open_test.go`
- Create: `internal/platform/river/client.go`
- Create: `internal/platform/river/transaction_test.go`
- Create: `apps/migrate/main.go`
- Modify: `go.mod`
- Modify: `go.sum`
- Test: `internal/platform/db/open_test.go`, `internal/platform/river/transaction_test.go`

- [ ] **Step 1: Write failing PostgreSQL integration tests.**

  Gate integration tests behind `TEST_DATABASE_URL`; skip only when that value
  is absent. The River test must create and clean up a test-only
  `transaction_probe` table with a setup connection, then begin one pgx-backed
  Ent transaction, insert a River job through River's transaction-aware API,
  insert a probe row, and roll back. A fresh connection must observe neither
  row. Repeat with commit and assert both the probe row and one `river_job` row
  exist. This proves the selected driver shares one physical transaction.

- [ ] **Step 2: Run the integration test with the required database URL.**

  Run: `rtk env TEST_DATABASE_URL="$DATABASE_URL" go test ./internal/platform/river -run TestEntAndRiverShareTransaction -count=1`  
  Expected: FAIL because Ent, River migrations, and the client factory do not
  exist.

- [ ] **Step 3: Add Ent source and deterministic generation.**

  Configure `entc.go` to generate into `internal/platform/db`, enable the
  required SQL/transaction feature flags, and start with an empty schema package
  so the baseline has no product model. Add a Taskfile command in a later task
  that runs exactly:

  ```bash
  rtk go run entgo.io/ent/cmd/ent generate --target internal/platform/db ./database/schema
  ```

  Keep `database/schema` as source and commit the resulting generated client.

- [ ] **Step 4: Implement migrations and the Go migration program.**

  `00001_river.sql` is generated from the pinned River version's official Goose
  migration support; copy the exact vendor migration sequence rather than
  hand-writing River tables. `00002_ent_schema.sql` is the empty Ent baseline
  migration and contains no product tables. `apps/migrate/main.go` loads only
  `database` config and invokes `goose.UpContext` against
  `database/migrations`; it never calls `ent.Schema.Create`.

- [ ] **Step 5: Implement one pgx pool for Ent and River.**

  `internal/platform/db/open.go` opens `*sql.DB` using `github.com/jackc/pgx/v5/stdlib`,
  configures bounded connection lifetime/idle settings from typed config, and
  constructs the generated Ent client from that driver. `internal/platform/river/client.go`
  constructs River with the same `*pgxpool.Pool` or compatible driver selected
  by the integration test; do not create separate uncoordinated pools in a
  transaction path. Expose an explicit transaction helper that receives both
  `*ent.Tx` and the matching River transaction executor.

- [ ] **Step 6: Run database and transaction verification.**

  Run:

  ```bash
  rtk go generate ./database/schema
  rtk gofmt -w apps/migrate database/schema internal/platform
  rtk go test ./internal/platform/db ./internal/platform/river -count=1
  rtk env TEST_DATABASE_URL="$DATABASE_URL" go test ./internal/platform/river -run TestEntAndRiverShareTransaction -count=1
  rtk go vet ./...
  ```

  Expected: committed generated Ent code is current, unit tests pass, and the
  rollback/commit integration assertions pass against PostgreSQL.

- [ ] **Step 7: Commit persistence infrastructure.**

  ```bash
  git add go.mod go.sum database apps/migrate internal/platform/db internal/platform/river
  git commit -m "feat: add ent goose and river persistence"
  ```

### Task 4: Create the Go mutation boundary and Huma v1 HTTP contract

**Files:**
- Create: `internal/usecase/transaction.go`
- Create: `internal/usecase/transaction_test.go`
- Create: `internal/transport/http/envelope.go`
- Create: `internal/transport/http/method/router.go`
- Create: `internal/transport/http/method/router_test.go`
- Create: `internal/transport/http/v1/routes.go`
- Create: `internal/transport/http/v1/system/status.go`
- Create: `internal/transport/http/api.go`
- Create: `internal/transport/http/api_test.go`
- Create: `internal/transport/http/openapi.go`
- Modify: `apps/api/main.go`
- Test: `internal/transport/http/**/*.go`, `internal/usecase/transaction_test.go`

- [ ] **Step 1: Write route-versioning tests before registration code.**

  Cover all three required path outcomes:

  ```go
  func TestV1GetBuildsAPIVersionedPath(t *testing.T) { /* GET /api/v1/status */ }
  func TestV1GetRejectsAPIPrefixedPath(t *testing.T) { /* /api/status => error */ }
  func TestV1GetRejectsVersionPrefixedPath(t *testing.T) { /* /v1/status => error */ }
  ```

  Add API tests for `GET /health`, `GET /health/ready`, and
  `GET /api/v1/status`. Assert status uses the exact envelope
  `{ "success": true, "data": { "status": "ok" } }`, while readiness returns
  HTTP 503 and `{ "success": false, "error": "NOT_READY", ... }` when its
  database probe fails.

- [ ] **Step 2: Run the focused tests to establish the red state.**

  Run: `rtk go test ./internal/usecase ./internal/transport/http/... -count=1`  
  Expected: FAIL because the usecase and Huma packages do not exist.

- [ ] **Step 3: Implement the only write-side transaction entrypoint.**

  Define the reusable dependency and callback shape in
  `internal/usecase/transaction.go`:

  ```go
  type Transaction func(context.Context, *ent.Tx, riverdriver.Executor) error
  type Runner interface { WithinTransaction(context.Context, Transaction) error }
  ```

  The implementation begins Ent and River work on the already-proven shared
  transaction, rolls back on callback error, commits only after callback
  success, and returns wrapped errors. Future intent-named usecases depend on
  `Runner`; no handler receives a write-capable Ent client.

- [ ] **Step 4: Implement Huma API construction and safe route wrapper.**

  `method.Router` stores a Huma API plus `Version("v1")`. `Get`, `Post`,
  `Patch`, `Put`, and `Delete` require a relative path beginning with `/`; they
  reject a path containing `/api` or any `v[0-9]+` segment and construct
  `"/api/" + version + relativePath`. Each registration also requires an
  explicit operation ID. `v1/routes.go` creates the v1 router and registers the
  `system` bundle. No feature package imports the raw Huma API.

  Keep `/health` and `/health/ready` outside the version wrapper. The API
  constructor installs request-ID middleware, structured request logging,
  standard error envelopes, and Huma's OpenAPI registration. Documentation is
  served at `/api/docs`; JSON is served at `/api/openapi.json`.

- [ ] **Step 5: Wire the API process without a TypeScript adapter.**

  `apps/api/main.go` loads `http_api`, builds the bootstrap runtime, creates the
  Huma API, starts `http.Server`, and performs graceful shutdown with a bounded
  context. It listens on `http_api.host:http_api.port`. Do not add CORS for the
  same-origin browser path; deployment edge configuration owns cross-origin
  policy if a later product requires it.

- [ ] **Step 6: Export a deterministic OpenAPI builder and verify the contract.**

  `internal/transport/http/openapi.go` writes the in-memory Huma specification
  to a supplied path with stable JSON formatting. The API test must parse this
  output and assert that `GET /api/v1/status` has operation ID
  `v1_get_system_status`, and that no `/api/status` path exists.

- [ ] **Step 7: Run focused API verification.**

  Run:

  ```bash
  rtk gofmt -w apps/api internal/usecase internal/transport/http
  rtk go test ./internal/usecase ./internal/transport/http/... -count=1
  rtk go vet ./...
  ```

  Expected: route guard, envelope, health, status, and generated OpenAPI tests
  pass.

- [ ] **Step 8: Commit the Go HTTP boundary.**

  ```bash
  git add apps/api internal/usecase internal/transport/http
  git commit -m "feat: add versioned huma api boundary"
  ```

### Task 5: Generate the OpenAPI-driven TypeScript client and proxy to Go

**Files:**
- Create: `tools/openapi/main.go`
- Create: `apps/web/hey-api.config.ts`
- Create: `apps/web/src/lib/api/generated/` (generator output)
- Create: `apps/web/src/lib/api/client.test.ts`
- Modify: `apps/web/src/lib/api/client.ts`
- Modify: `apps/web/src/lib/api/server.ts`
- Modify: `apps/web/src/server.ts`
- Delete: `apps/web/src/server/elysia-adapter.ts`
- Delete: `apps/web/src/server/elysia-adapter.test.ts`
- Modify: `apps/web/vite.config.ts`
- Modify: `apps/web/vite.config.test.ts`
- Modify: `apps/web/package.json`
- Modify: `package.json`
- Modify: `Taskfile.yml`
- Create: `contracts/openapi/.gitkeep`
- Test: `apps/web/src/lib/api/client.test.ts`, `apps/web/vite.config.test.ts`

- [ ] **Step 1: Write web tests for the generated-client wrapper and development proxy.**

  The client test creates the wrapper with `http://localhost:4101` and asserts
  that request URLs preserve `/api/v1/status`; it must not import an Elysia type
  or `@elysia/eden`. The Vite config test asserts a development proxy entry for
  both `/api` and `/health` with target `http://localhost:4101` and
  `changeOrigin: true`.

- [ ] **Step 2: Run the tests to show the pre-migration client is incompatible.**

  Run: `rtk bun test apps/web/src/lib/api/client.test.ts apps/web/vite.config.test.ts`  
  Expected: FAIL because the current client is Eden and Vite has no Go API proxy.

- [ ] **Step 3: Implement deterministic OpenAPI export and Hey API generation.**

  `tools/openapi/main.go` creates the same API object as `apps/api` without
  listening and writes `contracts/openapi/openapi.json`. Configure Hey API with
  the exact installed version and these plugins:

  ```ts
  import { defineConfig } from "@hey-api/openapi-ts";

  export default defineConfig({
    input: "../../contracts/openapi/openapi.json",
    output: { path: "src/lib/api/generated", format: "prettier" },
    plugins: ["@hey-api/client-fetch", "@hey-api/sdk", "@tanstack/react-query"],
  });
  ```

  Pin `@hey-api/openapi-ts` and the three named plugins to the same supported
  release line in `apps/web/package.json`. Commit the generated fetch client,
  SDK, types, and TanStack Query query-options/hooks.

- [ ] **Step 4: Replace Eden and the embedded server path.**

  `apps/web/src/lib/api/client.ts` configures the generated fetch client with
  the browser's same-origin base URL. `server.ts` configures a request-scoped
  generated client from `web.public.api_base_url`. `apps/web/src/server.ts`
  becomes the plain TanStack Start server entry; it does not inspect `/api` or
  import an API package. Remove `@repo/api`, `@elysia/eden`, Elysia, Prisma, and
  the adapter from web dependencies. Add the Vite development proxy.

- [ ] **Step 5: Add the three Taskfile contract commands.**

  Implement these exact command contracts:

  ```yaml
  api:openapi:
    cmds: ["rtk go run ./tools/openapi"]
  api:client:generate:
    cmds: ["rtk task api:openapi", "rtk bun --cwd apps/web run api:generate"]
  api:client:check:
    cmds: ["rtk task api:client:generate", "rtk git diff --exit-code -- contracts/openapi apps/web/src/lib/api/generated"]
  ```

  The `apps/web` script runs the pinned Hey API CLI. Do not fetch a live API
  server for generation.

- [ ] **Step 6: Run generation and web verification.**

  Run:

  ```bash
  rtk task api:client:generate
  rtk bun test apps/web/src/lib/api/client.test.ts apps/web/vite.config.test.ts
  rtk bun --cwd apps/web run check-types
  rtk task api:client:check
  ```

  Expected: OpenAPI and generated client are current; no Eden references remain
  in `apps/web/src`; web typecheck passes.

- [ ] **Step 7: Commit generated client integration.**

  ```bash
  git add tools/openapi contracts/openapi apps/web Taskfile.yml package.json bun.lock
  git rm apps/web/src/server/elysia-adapter.ts apps/web/src/server/elysia-adapter.test.ts
  git commit -m "feat: generate web api client from huma openapi"
  ```

### Task 6: Make realtime a versioned AsyncAPI boundary with Go publisher

**Files:**
- Create: `contracts/asyncapi/realtime.v1.yaml`
- Create: `internal/platform/realtime/publisher.go`
- Create: `internal/platform/realtime/publisher_test.go`
- Create: `internal/worker/river/realtime_publish.go`
- Create: `internal/worker/river/realtime_publish_test.go`
- Modify: `packages/realtime/src/events.ts`
- Modify: `packages/realtime/src/events.test.ts`
- Modify: `packages/realtime/src/index.ts`
- Delete: `packages/realtime/src/publisher.ts`
- Delete: `packages/realtime/src/publisher.test.ts`
- Modify: `apps/realtime/src/server.ts`
- Modify: `apps/realtime/src/server.test.ts`
- Modify: `apps/realtime/src/main.ts`
- Test: `internal/platform/realtime/publisher_test.go`, `internal/worker/river/realtime_publish_test.go`, `packages/realtime/src/events.test.ts`, `apps/realtime/src/server.test.ts`

- [ ] **Step 1: Define the failing cross-runtime event tests.**

  Define exactly one initial event type, `resource.updated.v1`, with
  `event_id`, `occurred_at`, `resource_id`, and `workspace_id`. The Go publisher
  test must assert `POST /internal/events`, JSON content type, the exact
  `x-realtime-api-key`, and an error for non-2xx responses. Node tests must
  reject the old unversioned `resource.updated` type and emit the validated v1
  event to both `resource:<resource_id>` and `workspace:<workspace_id>` rooms.

- [ ] **Step 2: Run the relevant tests to establish failure.**

  Run:

  ```bash
  rtk go test ./internal/platform/realtime ./internal/worker/river -count=1
  rtk bun test packages/realtime apps/realtime
  ```

  Expected: Go packages are absent and Node tests fail because the existing
  event is unversioned.

- [ ] **Step 3: Write AsyncAPI and aligned validators.**

  `contracts/asyncapi/realtime.v1.yaml` defines the private Go publisher channel
  `/internal/events`, HTTP `POST` operation, API-key header security scheme,
  and `resource.updated.v1` payload. Model the same required fields in Go and
  Zod. Keep the names in the wire payload snake_case; only local language types
  may use idiomatic casing.

- [ ] **Step 4: Implement the publisher and River delivery job.**

  The Go publisher API is:

  ```go
  type Event struct { Type, EventID, OccurredAt, ResourceID, WorkspaceID string }
  type Publisher interface { Publish(context.Context, Event) error }
  ```

  It uses a configured timeout and never exposes the internal key to the web.
  `realtime.publish.v1` is a typed River job whose handler validates the event
  and calls `Publisher.Publish`. A non-2xx response returns an error so River
  retries according to the job policy. The job itself performs no database
  mutation and does not open an Ent transaction.

- [ ] **Step 5: Preserve the Socket.IO process boundary.**

  Delete the now-unused TypeScript publisher; only the Go platform publisher
  may call the internal endpoint. Keep the Node `/internal/events` endpoint
  private at the edge. Update it to
  validate the versioned payload and emit the same payload under the
  `realtime-event` Socket.IO event name. Keep ticket authentication and product
  workspace authorization injectable; do not implement a baseline authorization
  rule beyond the existing deny-by-default behavior.

- [ ] **Step 6: Run realtime verification.**

  Run:

  ```bash
  rtk gofmt -w internal/platform/realtime internal/worker/river
  rtk go test ./internal/platform/realtime ./internal/worker/river -count=1
  rtk bun test packages/realtime apps/realtime
  rtk bun --cwd apps/realtime run check-types
  ```

  Expected: Go and TypeScript reject invalid/unversioned payloads, accept the
  same v1 contract, and a failed internal HTTP publish produces a retryable job
  handler error.

- [ ] **Step 7: Commit the realtime contract.**

  ```bash
  git add contracts/asyncapi internal/platform/realtime internal/worker/river packages/realtime apps/realtime
  git commit -m "feat: add asyncapi realtime publisher boundary"
  ```

### Task 7: Move queue execution and schedules into the Go River worker

**Files:**
- Create: `internal/worker/river/client.go`
- Create: `internal/worker/river/jobs.go`
- Create: `internal/worker/river/jobs_test.go`
- Create: `internal/worker/river/periodic.go`
- Create: `internal/worker/river/periodic_test.go`
- Modify: `apps/worker/main.go`
- Modify: `internal/bootstrap/bootstrap.go`
- Modify: `internal/bootstrap/bootstrap_test.go`
- Delete: `apps/scheduler/`
- Delete: `packages/queue/`
- Test: `internal/worker/river/jobs_test.go`, `internal/worker/river/periodic_test.go`

- [ ] **Step 1: Write River worker lifecycle and periodic-registration tests.**

  Test that `RegisterWorkers` installs `realtime.publish.v1` exactly once and
  test that `RegisterPeriodicJobs` receives the exact deterministic schedule
  set on every worker replica. The baseline schedule set is empty, so the test
  asserts an empty collection and establishes the registration extension point;
  it must not create a placeholder business schedule. Test that an unknown job
  name cannot be registered.

- [ ] **Step 2: Run the focused worker tests in the red state.**

  Run: `rtk go test ./internal/worker/river -run 'Test(RegisterWorkers|RegisterPeriodicJobs)' -count=1`  
  Expected: FAIL until the lifecycle and registration functions are added.

- [ ] **Step 3: Implement worker composition and graceful shutdown.**

  `apps/worker/main.go` loads `worker` configuration, starts River, calls
  `RegisterWorkers` and `RegisterPeriodicJobs`, waits for `SIGINT`/`SIGTERM`,
  then calls River's graceful stop with a bounded shutdown context before
  closing the database pool. It imports no HTTP handler and no product rule.

  `internal/worker/river/jobs.go` exposes only typed job registration. Future
  product handlers must decode arguments, issue any direct Ent reads they need,
  then invoke an `internal/usecase` mutation for a write. Document this rule
  adjacent to the registration API and test it through the transaction runner
  in Task 3 when the first product job is introduced.

- [ ] **Step 4: Implement periodic reconciliation semantics.**

  `periodic.go` registers schedules with River's leader-elected periodic API.
  Its API accepts only a job constructor and a stable cron or interval. Add a
  Go doc comment stating that any deadline-sensitive product task must enqueue
  a reconciliation job which scans due rows idempotently; a periodic tick may
  not be used as the sole source of truth. Do not add a scheduler process.

- [ ] **Step 5: Remove pg-boss and the TypeScript scheduler.**

  Delete `apps/scheduler` and `packages/queue`, remove their workspaces,
  scripts, Taskfile targets, Dockerfile, compose service/profile, ESLint files,
  and dependencies. Preserve no compatibility queue adapter: Go is now the
  only queue producer and consumer.

- [ ] **Step 6: Run worker verification.**

  Run:

  ```bash
  rtk gofmt -w apps/worker internal/bootstrap internal/worker/river
  rtk go test ./internal/bootstrap ./internal/worker/river -count=1
  rtk go vet ./...
  rtk rg -n "pg-boss|@repo/queue|apps/scheduler|@repo/scheduler" apps packages Taskfile.yml package.json docker-compose.yml Dockerfile.*
  ```

  Expected: Go tests and vet pass; the final search returns no active source or
  manifest references.

- [ ] **Step 7: Commit Go-only jobs and scheduling.**

  ```bash
  git add apps/worker internal/bootstrap internal/worker package.json Taskfile.yml docker-compose.yml bun.lock
  git rm -r apps/scheduler packages/queue Dockerfile.scheduler
  git commit -m "feat: move jobs and scheduling to go river worker"
  ```

### Task 8: Remove the TypeScript API/Prisma stack and make containers Go-native

**Files:**
- Create: `Dockerfile.migrate`
- Replace: `Dockerfile.api`
- Replace: `Dockerfile.worker`
- Modify: `Dockerfile.web`
- Modify: `docker-compose.yml`
- Modify: `Taskfile.yml`
- Modify: `package.json`
- Modify: `turbo.json`
- Modify: `apps/web/package.json`
- Delete: `apps/api/eslint.config.js`
- Delete: `apps/api/package.json`
- Delete: `apps/api/src/`
- Delete: `packages/application/`
- Delete: `packages/database/`
- Delete: `packages/storage/`
- Delete: `Dockerfile.scheduler`
- Create: `scripts/check-task-surface.ts`
- Test: Docker builds and `task --list`

- [ ] **Step 1: Write command-surface assertions.**

  Add `scripts/check-task-surface.ts`, which invokes `task --list` and checks
  it contains `dev:api`, `dev:worker`,
  `dev:migrate`, `api:openapi`, `api:client:generate`, `api:client:check`,
  `db:generate`, `db:migrate`, and `db:migrate:status`; it must not list
  `dev:scheduler`, `start:scheduler`, or Prisma commands.

- [ ] **Step 2: Run it and establish that existing tasks are obsolete.**

  Run: `rtk bun run scripts/check-task-surface.ts`  
  Expected: FAIL because the current task list still includes scheduler and
  Prisma task names and lacks the Go migration commands.

- [ ] **Step 3: Replace Docker images with minimal Go build stages.**

  `Dockerfile.api`, `Dockerfile.worker`, and `Dockerfile.migrate` use a Go
  builder stage to compile `./apps/api`, `./apps/worker`, and `./apps/migrate`,
  then copy only each static or CGO-compatible binary and `config/`/
  `database/migrations/` into a non-root runtime image. Choose one base image
  compatible with the selected pgx build mode and use it consistently. The web
  image no longer receives `DATABASE_URL`, generates Prisma, or copies backend
  packages; it loads only the public web configuration.

- [ ] **Step 4: Rewrite Compose and Taskfile around separate Go processes.**

  Compose starts `db`, one-shot `migrate`, `api`, `web`, optional `worker`, and
  optional `realtime`. `api` and `worker` depend on successful migration; web
  depends on API health, not on database. Expose API port 4101 for local
  development and keep the Socket.IO internal endpoint unexposed except through
  its intended local port/profile.

  Taskfile must use `go run ./apps/api`, `go run ./apps/worker`, and
  `go run ./apps/migrate` for development; `go build` for Go build targets;
  Goose commands for migrations; and the generation commands from Task 5.
  Remove all Bun/Turbo API, worker, scheduler, Prisma, and pg-boss commands.

- [ ] **Step 5: Delete obsolete application layers and dependency edges.**

  Remove the TypeScript files under `apps/api/src` while retaining Go
  `apps/api/main.go`. Remove `packages/application`, `packages/database`, and
  `packages/storage` after the web no longer imports them. Remove Elysia, Eden,
  Prisma, pg-boss, storage package dependencies, and their transitive-only
  configuration from manifests and regenerate `bun.lock`. Do not delete
  `apps/realtime` or `packages/realtime`.

- [ ] **Step 6: Verify images and command surface.**

  Run:

  ```bash
  rtk task --list
  rtk task db:generate
  rtk task api:client:check
  rtk docker compose build migrate api worker web realtime
  rtk rg -n "Elysia|@elysia|Prisma|prisma|pg-boss|Eden|@repo/api|@repo/database|@repo/application|@repo/storage" apps packages Dockerfile.* docker-compose.yml Taskfile.yml package.json turbo.json
  ```

  Expected: required commands exist; all listed images build; the final search
  returns no active architecture reference except intentional migration-history
  comments, which must be removed rather than suppressed.

- [ ] **Step 7: Commit runtime cleanup.**

  ```bash
  git add Dockerfile.api Dockerfile.worker Dockerfile.migrate Dockerfile.web docker-compose.yml Taskfile.yml package.json turbo.json apps/api/main.go apps/web bun.lock
  git rm -r apps/api/src packages/application packages/database packages/storage Dockerfile.scheduler
  git rm apps/api/eslint.config.js apps/api/package.json apps/api/tsconfig.json
  git commit -m "refactor: replace typescript backend runtimes with go"
  ```

### Task 9: Align documentation, agent rules, and full verification

**Files:**
- Create: `scripts/check-architecture-docs.ts`
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `.agent/architecture.md`
- Modify: `.agent/api.md`
- Modify: `.agent/database.md`
- Modify: `.agent/config.md`
- Modify: `.agent/worker.md`
- Delete: `.agent/scheduler.md`
- Test: repository quality commands and fresh Compose smoke test

- [ ] **Step 1: Write documentation consistency checks.**

  Add a repository documentation test or scripted check that requires README
  and agent rules to mention Huma, Ent, Goose, River, `/api/v1`, Hey API, and
  Socket.IO; it must reject the active-architecture terms Elysia, Prisma,
  pg-boss, Eden, and standalone scheduler. Exclude dated historical specs from
  this check by path, not by a broad textual exception.

- [ ] **Step 2: Run the check before rewriting documentation.**

  Run: `rtk bun run scripts/check-architecture-docs.ts`  
  Expected: FAIL until README and `.agent` rules describe the new architecture.

- [ ] **Step 3: Rewrite contributor-facing architecture guidance.**

  Update README diagrams and command reference to the target runtimes. Update
  `AGENTS.md` and `.agent` files so they enforce: Go-only domain mutations;
  Ent direct reads; Goose-only schema changes; River scheduling in worker;
  private HTTP realtime publishing; `/api/v1` through the route wrapper; and
  Huma-to-Hey generation. Delete the scheduler-specific rule.

- [ ] **Step 4: Run all static and focused checks.**

  Run:

  ```bash
  rtk go test ./... -count=1
  rtk go vet ./...
  rtk bun test
  rtk bun run lint
  rtk bun run check-types
  rtk task api:client:check
  rtk git diff --check
  ```

  Expected: each command exits 0. If a Bun test refers to a removed runtime,
  delete or replace that test in the owning migration task rather than skipping
  it here.

- [ ] **Step 5: Run a clean PostgreSQL Compose smoke test.**

  Run:

  ```bash
  rtk docker compose down --volumes
  rtk docker compose up --build -d migrate api web
  rtk curl --fail http://localhost:4101/health
  rtk curl --fail http://localhost:4101/health/ready
  rtk curl --fail http://localhost:4101/api/v1/status
  rtk task api:client:check
  rtk docker compose down --volumes
  ```

  Expected: a fresh database is migrated by Goose (including River tables),
  Go health/readiness/status return success, and client regeneration has no
  diff. The final down command removes only the plan's test Compose volume.

- [ ] **Step 6: Commit documentation and final migration verification.**

  ```bash
  git add README.md AGENTS.md .agent config/README.md Taskfile.yml
  git rm .agent/scheduler.md
  git commit -m "docs: document go river backend architecture"
  ```

## Final acceptance checklist

- [ ] `go test ./... -count=1` and `go vet ./...` pass.
- [ ] `bun test`, `bun run lint`, and `bun run check-types` pass for retained TypeScript workspaces.
- [ ] `task api:openapi`, `task api:client:generate`, and `task api:client:check` pass with no diff.
- [ ] The Ent/River PostgreSQL integration test proves commit and rollback atomicity.
- [ ] A fresh Compose database is migrated only by Goose and contains River tables.
- [ ] `GET /api/v1/status`, `/health`, and `/health/ready` pass smoke tests; no `/api/status` route exists.
- [ ] No active source, manifest, container, config, or task references Elysia, Prisma, pg-boss, Eden, Redis/Asynq, or `apps/scheduler`.
- [ ] The only cross-language asynchronous contract is the versioned AsyncAPI realtime event; River jobs are Go-only.

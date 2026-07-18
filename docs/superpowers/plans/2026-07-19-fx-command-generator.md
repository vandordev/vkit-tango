# Fx Command and Generator Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace manual Go runtime wiring with Uber Fx, reusable commands, Chi/Huma HTTP adapters, River job/scheduler roots, and Taskfile-driven `vx` scaffolding plus deterministic registries.

**Architecture:** `internal/contract` defines the command and adapter boundaries. Fx modules in each runtime root select shared infrastructure plus generated providers and registrations. `vpkg/vandor/go` templates create one focused surface at a time, while a Go AST scanner generates Fx, route, worker, and scheduler registries from those fixed boundaries.

**Tech Stack:** Go 1.25, `go.uber.org/fx`, Chi v5, Huma v2 with `humachi`, River, Ent, Goose, Taskfile, `vx`, Go AST/format.

---

## File Structure

```text
apps/
  api/main.go                         # Fx API process
  worker/main.go                      # Fx River worker process
  scheduler/main.go                   # new Fx periodic-enqueue process
  migrate/main.go                     # retained non-Fx migration process

internal/
  app/
    common.go                          # shared configuration/environment helpers
    api.go                             # API Fx module
    worker.go                          # worker Fx module
    scheduler.go                       # scheduler Fx module
  contract/
    command.go                         # generic command boundary
    http.go                            # route-registration boundary
    job.go                             # worker-registration boundary
    scheduler.go                       # periodic-registration boundary
  generated/fx/
    usecases_gen.go                    # generated Fx command providers
    http_gen.go                        # generated HTTP handler providers/invokes
    worker_gen.go                      # generated River job registration
    scheduler_gen.go                   # generated periodic registration
  transport/http/
    api.go                             # humachi API creation
    server.go                          # Chi router and lifecycle hook
    method/register.go                 # typed GET/POST/PUT/PATCH/DELETE helpers
    handler/system_metadata/
      set_system_metadata.go           # one PUT operation adapter
  usecase/
    set_system_metadata.go             # renamed command implementation
    set_system_metadata_test.go        # command behavior tests
  worker/river/
    system_metadata.go                 # command-backed job adapter
    realtime_publish.go                # realtime job adapter
  scheduler/river/
    module.go                          # generated registry consumes implementations

vpkg/vandor/go/
  vpkg.yaml
  templates/{usecase,http_handler,job,scheduler}.vxt
  tools/sync/main.go
  tools/sync/main_test.go
```

Generated files under `internal/generated/fx` are committed output and include
the standard generated-file banner. They are written only by `task sync:*`.

### Task 1: Add Uber Fx and shared contract boundaries

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`
- Create: `internal/contract/command.go`
- Create: `internal/contract/command_test.go`
- Create: `internal/contract/http.go`
- Create: `internal/contract/job.go`
- Create: `internal/contract/scheduler.go`

- [ ] **Step 1: Write the failing command-contract test**

Create `internal/contract/command_test.go`:

```go
package contract_test

import (
	"context"
	"testing"

	"github.com/vandordev/vkit-tango/internal/contract"
)

type testCommand struct{}

func (testCommand) Execute(context.Context, string) (int, error) { return 1, nil }

var _ contract.Command[string, int] = testCommand{}

func TestCommandExecutesTypedInput(t *testing.T) {
	got, err := testCommand{}.Execute(context.Background(), "input")
	if err != nil || got != 1 {
		t.Fatalf("Execute() = (%d, %v), want (1, nil)", got, err)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `rtk go test ./internal/contract -run TestCommandExecutesTypedInput -count=1`

Expected: FAIL because `internal/contract` does not exist.

- [ ] **Step 3: Define the four minimal contracts**

Create `internal/contract/command.go`:

```go
package contract

import "context"

type Command[I any, O any] interface {
	Execute(context.Context, I) (O, error)
}
```

Define `HTTPHandler` in `http.go` with `RegisterRoutes()`. Define
`WorkerRegistrar` in `job.go` with `RegisterWorkers(*riverqueue.Workers)`.
Define `SchedulerRegistrar` in `scheduler.go` with
`RegisterPeriodicJobs() []*riverqueue.PeriodicJob`. Import River only in the
two River-specific contract files; do not put HTTP or River dependencies into
`command.go`.

- [ ] **Step 4: Add Fx dependency and format module files**

Run: `rtk go get go.uber.org/fx@v1.24.0 && rtk gofmt -w internal/contract`

Expected: `go.mod` directly requires `go.uber.org/fx v1.24.0`; contract files
are formatted.

- [ ] **Step 5: Run contract and module checks**

Run: `rtk go test ./internal/contract -count=1 && rtk go mod tidy && rtk git diff --check`

Expected: PASS with no whitespace errors.

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum internal/contract
git commit -m "feat: add Fx command contracts"
```

### Task 2: Convert the baseline mutation into a command

**Files:**
- Modify: `internal/usecase/system_metadata.go`
- Modify: `internal/usecase/system_metadata_test.go`
- Modify: `internal/usecase/transaction.go`
- Create: `internal/usecase/module.go`
- Create: `internal/usecase/module_test.go`

- [ ] **Step 1: Write the failing constructor/contract test**

Add to `internal/usecase/system_metadata_test.go`:

```go
func TestNewSetSystemMetadataImplementsCommand(t *testing.T) {
	command := usecase.NewSetSystemMetadata(usecase.Runner{})
	var _ contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult] = command
}
```

Use package `usecase_test` and import `internal/contract` and
`internal/usecase`; this proves the public command surface without accessing
fields. Do not test the HTTP adapter here.

- [ ] **Step 2: Run test to verify it fails**

Run: `rtk go test ./internal/usecase -run TestNewSetSystemMetadataImplementsCommand -count=1`

Expected: FAIL because `NewSetSystemMetadata` does not exist.

- [ ] **Step 3: Replace the local interface/service name with one command struct**

Keep `SetSystemMetadataInput` and `SetSystemMetadataResult`. Delete the local
`SetSystemMetadata` interface and rename `SystemMetadataService` to
`SetSystemMetadata`. Add:

```go
func NewSetSystemMetadata(runner Runner) *SetSystemMetadata {
	return &SetSystemMetadata{runner: runner}
}
```

Keep the existing transaction, Ent upsert, and transactional realtime enqueue
semantics exactly intact: change the receiver from `service SystemMetadataService`
to `command *SetSystemMetadata`, replace every `service.Runner` reference
with `command.runner`, and preserve the existing return values and errors.
Add a compile assertion in the production file:

```go
var _ contract.Command[SetSystemMetadataInput, SetSystemMetadataResult] = (*SetSystemMetadata)(nil)
```

- [ ] **Step 4: Add the usecase Fx module and constructor for Runner**

In `transaction.go`, add `NewRunner(database *sql.DB, river *riverqueue.Client[*sql.Tx]) Runner` that returns the existing `Runner` value. In
`module.go`, define:

```go
var Module = fx.Options(
	fx.Provide(NewRunner),
)
```

Do not hand-register `NewSetSystemMetadata` here: that provider belongs in the
generated registry created in Task 6.

- [ ] **Step 5: Run focused behavior tests**

Run: `rtk go test ./internal/usecase ./internal/platform/river -count=1`

Expected: PASS; the existing transaction test still proves that Ent mutation
and River enqueue share a SQL transaction.

- [ ] **Step 6: Commit**

```bash
git add internal/usecase
git commit -m "refactor: model metadata mutation as command"
```

### Task 3: Replace the HTTP mux with Chi + Huma via humachi

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`
- Modify: `internal/transport/http/api.go`
- Modify: `internal/transport/http/api_test.go`
- Modify: `internal/transport/http/method/router.go`
- Modify: `internal/transport/http/method/router_test.go`
- Create: `internal/transport/http/method/register.go`
- Create: `internal/transport/http/method/register_test.go`
- Create: `internal/transport/http/server.go`
- Create: `internal/transport/http/module.go`
- Create: `internal/transport/http/handler/system_metadata/set_system_metadata.go`

- [ ] **Step 1: Write failing route helper tests**

In `method/register_test.go`, create a Chi router, wrap it with
`humachi.New(router, huma.DefaultConfig("test", "1.0.0"))`, register a typed
`PUT` route, and issue an `httptest` request. Assert HTTP 200 and that the
response body has the expected Huma output. Add a second test asserting that
the generated OpenAPI contains the method/path operation ID.

- [ ] **Step 2: Run test to verify it fails**

Run: `rtk go test ./internal/transport/http/method -run TestPUT -count=1`

Expected: FAIL because Chi, `humachi`, and `method.PUT` are unavailable.

- [ ] **Step 3: Add Chi and humachi dependencies and method helpers**

Run: `rtk go get github.com/go-chi/chi/v5 github.com/danielgtaylor/huma/v2/adapters/humachi`

Create typed `GET`, `POST`, `PUT`, `PATCH`, and `DELETE` helpers that all call
one unexported `register` function. That function sets an operation ID derived
from lower-case method plus normalized path, assigns method/path/summary/tags,
and invokes `huma.Register`.

- [ ] **Step 4: Implement the Chi/Huma transport module**

Create a Chi router with request ID, real-IP, recover, and timeout middleware.
Create `huma.API` with `humachi.New`, disable the Huma docs route, retain
`/api/openapi.json`, and expose the base Huma API without an implicit path group.
Register `/health` and `/health/ready` directly on Chi. Provide an Fx lifecycle-managed `http.Server`
whose `OnStart` calls `ListenAndServe` in a goroutine and whose `OnStop`
calls `Shutdown` with the configured timeout.

- [ ] **Step 5: Move metadata mutation to its one-operation handler**

Create `SetSystemMetadataHandler` with an injected
`contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult]`.
Its `RegisterRoutes` calls `method.PUT` for
`/api/v1/system-metadata/{key}` on the base Huma API; its `Handle` maps the
Huma input body to the command input and maps the result to the existing
success envelope. Keep all request/response structs unexported.

The generated HTTP registry will construct and invoke this handler in Task 6;
do not add a manually maintained route list.

- [ ] **Step 6: Update transport contract tests**

Replace direct calls to `transport.NewHandler` in `api_test.go` with an Fx
test app that provides a fake `contract.Command`, invokes the generated-style
handler registration, and serves the Chi router. Preserve assertions for
`/health`, `/health/ready`, `/api/v1/status`, the metadata `PUT`, and
`/api/openapi.json`.

- [ ] **Step 7: Run HTTP checks**

Run: `rtk go test ./internal/transport/http ./internal/transport/http/method ./internal/transport/http/handler/system_metadata -count=1 && rtk go test ./tools/openapi -count=1`

Expected: PASS; OpenAPI remains generated from Huma's API instance.

- [ ] **Step 8: Commit**

```bash
git add go.mod go.sum internal/transport/http
git commit -m "feat: serve Huma routes through Chi"
```

### Task 4: Make River adapters conform to shared contracts

**Files:**
- Modify: `internal/worker/river/system_metadata.go`
- Modify: `internal/worker/river/realtime_publish.go`
- Modify: `internal/worker/river/jobs_test.go`
- Modify: `internal/worker/river/realtime_publish_test.go`
- Modify: `internal/worker/river/periodic.go`
- Modify: `internal/worker/river/periodic_test.go`
- Create: `internal/scheduler/river/module.go`

- [ ] **Step 1: Write failing contract assertions for workers and schedules**

Add compile assertions in `jobs_test.go` that the metadata and realtime
registrars implement `contract.WorkerRegistrar`. Add an assertion in
`periodic_test.go` that the baseline scheduler registrar implements
`contract.SchedulerRegistrar` and returns an empty job list.

- [ ] **Step 2: Run test to verify it fails**

Run: `rtk go test ./internal/worker/river -run 'Test(Register|SetSystem)' -count=1`

Expected: FAIL because current variadic registration functions do not expose
the shared registrar contract.

- [ ] **Step 3: Replace variadic registration with concrete registrars**

Create a concrete metadata job registrar that holds
`contract.Command[SetSystemMetadataInput, SetSystemMetadataResult]`, adds its
typed worker to a supplied `*riverqueue.Workers`, and returns no error. Create
a concrete realtime registrar around `platformrealtime.Publisher` that adds
`RealtimePublishWorker`.

Keep `Work` methods small: metadata maps typed args to the shared command and
realtime calls only `Publisher.Publish`. Delete the variadic optional-command
API from `RegisterWorkers`.

- [ ] **Step 4: Move periodic ownership to the scheduler boundary**

Replace `RegisterPeriodicJobs()` with an exported baseline scheduler registrar
in `internal/scheduler/river`. It returns an empty slice now but implements
`contract.SchedulerRegistrar`; the generated scheduler registry is its only
future aggregation point. Remove periodic registration from the worker app
module.

- [ ] **Step 5: Run River adapter tests**

Run: `rtk go test ./internal/worker/river ./internal/scheduler/river ./internal/platform/river -count=1`

Expected: PASS, including the existing transaction integration coverage.

- [ ] **Step 6: Commit**

```bash
git add internal/worker/river internal/scheduler/river
git commit -m "refactor: register River adapters through contracts"
```

### Task 5: Create the `vpkg` scaffold package and preview tests

**Files:**
- Create: `vpkg/vandor/go/vpkg.yaml`
- Create: `vpkg/vandor/go/templates/usecase.vxt`
- Create: `vpkg/vandor/go/templates/http_handler.vxt`
- Create: `vpkg/vandor/go/templates/job.vxt`
- Create: `vpkg/vandor/go/templates/scheduler.vxt`
- Create: `vpkg/vandor/go/README.md`

- [ ] **Step 1: Write the manifest with explicit template exports**

Set the manifest name to `vandor/go`, kind to `template-pack`, and exports to
`usecase`, `http-handler`, `job`, and `scheduler`. Make `usecase` the
default only if it is unambiguous; otherwise require the explicit export in every
Taskfile call.

- [ ] **Step 2: Write the usecase template with its mandatory test file**

`usecase.vxt` accepts `name` and writes
`internal/usecase/{{ name | snake }}.go` and
`internal/usecase/{{ name | snake }}_test.go`. The production file defines
`{{ name | pascal }}Input`, `{{ name | pascal }}Result`, one public
`{{ name | pascal }}` struct, `New{{ name | pascal }}`, and a compile assertion
against `contract.Command`. The test file constructs that command with its
declared dependencies and asserts the intended command contract.

- [ ] **Step 3: Write the three thin adapter templates**

`http_handler.vxt` requires only `name`, `method`, and `path`. It writes one
handler struct whose name is `{{ name | pascal }}Handler`, injects the matching
generic command contract, and registers through
`method.{{ method | upper }}`. It writes no test file.

`job.vxt` requires `name`, writes one typed River worker/registrar that invokes
the matching command, and writes no test file. `scheduler.vxt` requires `name`,
writes one scheduler registrar that enqueues the corresponding job, and writes no
test file. Each template uses only paths rooted at `project.go.module_root`.

- [ ] **Step 4: Validate every template with non-destructive vx plans**

Run:

```bash
rtk vx view vandor/go:usecase --plan --set name=ExampleCommand
rtk vx view vandor/go:http-handler --plan --set name=ExampleCommand --set method=PUT --set path=/api/v1/examples/{id}
rtk vx view vandor/go:job --plan --set name=ExampleCommand
rtk vx view vandor/go:scheduler --plan --set name=ExampleCommand
```

Expected: each command prints only planned paths; no project file is written.

- [ ] **Step 5: Document the Taskfile-only workflow**

In `vpkg/vandor/go/README.md`, document that package maintainers may use `vx`
for preview, while application developers use `task add:*` and `task sync:*`.
List required values and failure behavior for existing generated paths.

- [ ] **Step 6: Commit**

```bash
git add vpkg/vandor/go
git commit -m "feat: add Go command scaffold package"
```

### Task 6: Implement deterministic Fx registry synchronization

**Files:**
- Create: `vpkg/vandor/go/tools/sync/main.go`
- Create: `vpkg/vandor/go/tools/sync/main_test.go`
- Create: `vpkg/vandor/go/tools/sync/testdata/project/internal/usecase/example.go`
- Create: `vpkg/vandor/go/tools/sync/testdata/project/internal/transport/http/handler/example.go`
- Create: `vpkg/vandor/go/tools/sync/testdata/project/internal/worker/river/example.go`
- Create: `vpkg/vandor/go/tools/sync/testdata/project/internal/scheduler/river/example.go`
- Create: `internal/generated/fx/usecases_gen.go`
- Create: `internal/generated/fx/http_gen.go`
- Create: `internal/generated/fx/worker_gen.go`
- Create: `internal/generated/fx/scheduler_gen.go`

- [ ] **Step 1: Write fixture-driven failing generator tests**

Use a temporary copied fixture project containing one valid use case, one HTTP
handler, one worker registrar, and one scheduler registrar. Test that
`sync --surface usecase` emits sorted `fx.Provide` entries with `fx.As`
command bindings; `http` emits both provider and `RegisterRoutes` invoke;
`worker` emits worker registrar invoke; and `scheduler` emits periodic
registrar invoke.

Add invalid fixtures for a missing `New<name>` constructor and a type missing its
required contract assertion. Assert a descriptive error naming the file and
missing contract.

- [ ] **Step 2: Run test to verify it fails**

Run: `rtk go test ./vpkg/vandor/go/tools/sync -count=1`

Expected: FAIL because no sync tool exists.

- [ ] **Step 3: Implement AST discovery and validation**

Parse only non-test `.go` files in the four documented surface directories.
For each file, find its exported `New<name>` function and the corresponding public
struct. Require a production compile assertion matching the surface contract;
do not infer membership from filename alone. Sort imports and providers by
import path and constructor name. Use `go/format` before atomically writing
each output file.

Support `--surface usecase|http|worker|scheduler|all`, defaulting to `all`.
Write exactly the selected registry files and create `internal/generated/fx`
when absent. Every generated file begins:

```go
// Code generated by vandor sync. DO NOT EDIT.
```

- [ ] **Step 4: Generate the baseline registries**

Run: `rtk go run ./vpkg/vandor/go/tools/sync --surface all`

Expected: four formatted registry files appear under `internal/generated/fx`;
the usecase registry contains `NewSetSystemMetadata`, the HTTP registry contains
`SetSystemMetadataHandler`, and worker/scheduler registries contain their
baseline registrars.

- [ ] **Step 5: Run generator and generated-package tests**

Run: `rtk go test ./vpkg/vandor/go/tools/sync ./internal/generated/fx -count=1 && rtk gofmt -w internal/generated/fx`

Expected: PASS with no diff from a second identical sync invocation.

- [ ] **Step 6: Commit**

```bash
git add vpkg/vandor/go/tools/sync internal/generated/fx
git commit -m "feat: generate Fx surface registries"
```

### Task 7: Build the Fx infrastructure and three composition roots

**Files:**
- Create: `internal/app/common.go`
- Create: `internal/app/api.go`
- Create: `internal/app/worker.go`
- Create: `internal/app/scheduler.go`
- Create: `internal/app/module_test.go`
- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`
- Create: `apps/scheduler/main.go`
- Create: `apps/scheduler/main_test.go`
- Modify: `apps/api/main.go`
- Modify: `apps/api/main_test.go`
- Modify: `apps/worker/main.go`
- Modify: `Taskfile.yml`

- [ ] **Step 1: Write failing Fx graph tests with replacement providers**

Create `internal/app/module_test.go`. For each root module, construct
`fx.New(api.Module, testInfrastructure(), fx.Invoke(func(httpcontract.API) {}))` with its module plus `fx.Replace` providers for database, Ent,
River producer/client, and HTTP server so no external connection is opened.
Invoke a sentinel that requests the required root dependency:

```go
func TestAPIModuleBuilds(t *testing.T) {
	app := fx.New(APIModule, testInfrastructure(), fx.Invoke(func(httpcontract.API) {}))
	if err := app.Err(); err != nil {
		t.Fatal(err)
	}
}
```

Define equivalent `TestWorkerModuleBuilds` and `TestSchedulerModuleBuilds`.
`testInfrastructure()` must use `fx.Replace` values of the exact concrete
types used by the real providers. Do not start the app in these graph tests.

- [ ] **Step 2: Run test to verify it fails**

Run: `rtk go test ./internal/app -run 'Test(API|Worker|Scheduler)ModuleBuilds' -count=1`

Expected: FAIL because `internal/app` and the root modules do not exist.

- [ ] **Step 3: Add scheduler configuration and common Fx providers**

Add `config.Scheduler` and `LoadScheduler`, loading `app`, `database`,
`worker`, `realtime`, and `observability` so the scheduler uses the same
database and River settings as the worker. Extend configuration tests with a
fixture asserting that `max_workers` is parsed and available.

In `internal/app/common.go`, provide environment capture, typed config loading,
`postgres.Open` with an Fx lifecycle close hook for the Ent client, a River
producer, and `usecase.Module`. Each provider must return an error rather than
calling `log.Fatal` or `panic`.

- [ ] **Step 4: Define the API, worker, and scheduler Fx modules**

`internal/app/api.go` exports `APIModule`, which composes common providers,
`generatedfx.UsecaseModule`, the HTTP transport module, and `generatedfx.HTTPModule`.

`internal/app/worker.go` exports `WorkerModule`, which composes common providers, `generatedfx.UsecaseModule`,
the realtime publisher, the worker River client, and `generatedfx.WorkerModule`.
Append lifecycle hooks that start and stop the River worker client with the
configured shutdown timeout.

`internal/app/scheduler.go` exports `SchedulerModule`, which composes common providers, an empty River worker
set, the River periodic client, and `generatedfx.SchedulerModule`. Its lifecycle
starts/stops the periodic client but never registers executable workers.

- [ ] **Step 5: Replace long-lived mains with Fx roots**

Each main must only call its root:

```go
// apps/api/main.go
func main() { fx.New(app.APIModule).Run() }

// apps/worker/main.go
func main() { fx.New(app.WorkerModule).Run() }

// apps/scheduler/main.go
func main() { fx.New(app.SchedulerModule).Run() }
```

Use aliases so each app imports its own module unambiguously. Remove duplicated
signal handling, environment copying, manual `Close`, `Start`, and `Stop`
code from API and worker mains. Keep `apps/migrate/main.go` unchanged.

- [ ] **Step 6: Add scheduler Taskfile and build wiring**

Add `dev:scheduler` using `rtk go run ./apps/scheduler`. Add
`./apps/scheduler` to the Go build command. Do not alter web/realtime tasks.

- [ ] **Step 7: Run module and entrypoint checks**

Run: `rtk go test ./internal/app ./apps/api ./apps/worker ./apps/scheduler ./internal/config -count=1 && rtk go build ./apps/api ./apps/worker ./apps/scheduler`

Expected: PASS without a database connection from graph tests.

- [ ] **Step 8: Commit**

```bash
git add apps/api apps/worker apps/scheduler internal/app internal/config Taskfile.yml
git commit -m "feat: add Fx runtime composition roots"
```

### Task 8: Expose all generation and sync operations through Taskfile

**Files:**
- Modify: `Taskfile.yml`
- Modify: `README.md`
- Create: `scripts/check-generated.sh`

- [ ] **Step 1: Add explicit Taskfile input validation**

Add `add:usecase`, `add:http-handler`, `add:job`, and `add:scheduler` tasks.
Each requires its documented variables and invokes one explicit target: `rtk vx gen vandor/go:usecase --apply`, `rtk vx gen vandor/go:http-handler --apply`, `rtk vx gen vandor/go:job --apply`, or `rtk vx gen vandor/go:scheduler --apply`, each with explicit `--set` flags. `add:http-handler` validates `method` against
`GET|POST|PUT|PATCH|DELETE` before invoking `vx` and passes `name`, `method`,
and `path` only; it never accepts a `command` variable.

- [ ] **Step 2: Add surface and umbrella sync tasks**

Define `sync:usecase`, `sync:http`, `sync:worker`, and `sync:scheduler` as
`rtk go run ./vpkg/vandor/go/tools/sync --surface <surface>`. Define `sync` as
their ordered dependency sequence followed by `rtk gofmt -w internal/generated/fx`.
Each `add:*` runs exactly its related `sync:*` after a successful `vx` apply.

- [ ] **Step 3: Write the generated-output check script**

Create `scripts/check-generated.sh` to run `task sync` and then
`git diff --exit-code -- internal/generated/fx`. The script exits non-zero with
the generated diff when committed output is stale. Make it executable.

- [ ] **Step 4: Exercise safe Taskfile paths**

Run:

```bash
rtk task sync
rtk ./scripts/check-generated.sh
rtk task add:http-handler method=TRACE path=/api/v1/test name=Example
```

Expected: the first two commands pass without a diff; the final command fails
before calling `vx` with an unsupported-method error and writes no file.

- [ ] **Step 5: Document the public workflow**

Update `README.md` commands and architecture sections with the four add tasks,
the four targeted sync tasks, `task sync`, the three Fx roots, and the rule
that only use cases require a paired test file. Do not document direct `vx`
commands as application-developer workflow.

- [ ] **Step 6: Commit**

```bash
git add Taskfile.yml README.md scripts/check-generated.sh
git commit -m "feat: add Taskfile command generation workflow"
```

### Task 9: Regenerate OpenAPI and perform end-to-end verification

**Files:**
- Modify: `contracts/openapi/openapi.json`
- Modify: `apps/api/main_test.go`
- Modify: `apps/worker/main_test.go` (create if absent)
- Modify: `apps/scheduler/main_test.go`

- [ ] **Step 1: Add root smoke tests that do not reach external services**

Use the same replacement-provider pattern from Task 7 to assert that each
entrypoint module can be constructed, including generated registries. API
asserts that the Huma API has `/api/v1/system-metadata/{key}`; worker asserts
that the metadata and realtime jobs register; scheduler asserts that the empty
baseline periodic schedule registers.

- [ ] **Step 2: Run the focused runtime checks**

Run: `rtk go test ./apps/api ./apps/worker ./apps/scheduler ./internal/app -count=1`

Expected: PASS without a PostgreSQL, River, or HTTP listener connection.

- [ ] **Step 3: Regenerate and verify OpenAPI**

Run: `rtk task api:openapi && rtk git diff --check -- contracts/openapi/openapi.json`

Expected: generated document keeps `/api/v1/system-metadata/{key}` with a PUT
operation and does not expose unversioned business endpoints.

- [ ] **Step 4: Run the complete repository gates**

Run: `rtk task sync && rtk task quality && rtk task build`

Expected: all Go/TypeScript checks, generated-client freshness, and all runtime
builds pass. Inspect `git status --short` and ensure remaining changes are the
expected source, generated registries, OpenAPI artifact, and documentation.

- [ ] **Step 5: Commit final generated artifacts and verification changes**

```bash
git add apps contracts/openapi internal/generated/fx
git commit -m "test: verify Fx command runtime"
```

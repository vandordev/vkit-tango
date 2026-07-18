# Agent Development Guidance Implementation Plan

> **For agentic workers:** Execute this plan inline with
> `superpowers:executing-plans`; do not spawn subagents or create a worktree.
> Keep `apps/web/src/routeTree.gen.ts` outside this plan because it is already
> modified before implementation begins.

**Goal:** Give agents enforceable, focused guidance for architecture, typed
contracts, Huma metadata, safe data changes, and verification without changing
product behaviour or inventing a schema for system metadata.

**Architecture:** Keep `AGENTS.md` as a concise index and add three focused
`.agent` references: repository ownership, change workflows, and verification.
Strengthen the shared Huma method helper so generated and handwritten handlers
derive one domain tag from a `/api/v1` path. Documentation checks assert that
the guidance remains discoverable and does not regress to the prior runtime
boundaries.

**Tech Stack:** Go, Huma/humachi, Uber Fx, Ent, Goose, River, Taskfile, vx,
Bun/TypeScript, `samber/lo` policy documentation.

---

## File Structure

- Modify: `AGENTS.md` — links to focused agent references and concise global
  guardrails.
- Create: `.agent/repository-map.md` — reference for ownership, source of
  truth, generated output, `internal/lib`, and `samber/lo`.
- Create: `.agent/workflows.md` — how-to guides for use cases, Huma, River,
  Ent/Goose, web contracts, local development, errors, and security.
- Create: `.agent/verification.md` — change-category verification matrix.
- Modify: `.agent/architecture.md`, `.agent/api.md`, `.agent/database.md`,
  `.agent/worker.md`, `.agent/web.md`, `.agent/config.md` — add only their
  specific new constraints and point to shared guidance where appropriate.
- Modify: `scripts/check-architecture-docs.ts` — assert guidance files and
  stable terms are present; reject superseded architecture phrases.
- Modify: `internal/transport/http/method/register.go` — derive and validate
  exactly one resource tag from an API v1 path.
- Modify: `internal/transport/http/method/register_test.go` — test derived
  tag and invalid path behaviour alongside deterministic operation IDs.
- Modify: `vpkg/vandor/go/templates/http_handler.vxt` — scaffold named HTTP
  DTOs and pass the deterministic tag helper instead of `nil`.
- Modify: `vpkg/vandor/go/README.md` — document the unchanged three-argument
  Taskfile interface and derived tag convention.
- Create: `scripts/check-vpkg-templates.ts` — assert the HTTP template keeps
  named DTOs and `method.Tags`; run a vx plan to prove it renders.
- Modify: `Taskfile.yml` — run the template check under `task test`.

## Task 1: Make the documentation guidance discoverable

**Files:**
- Create: `.agent/repository-map.md`
- Create: `.agent/workflows.md`
- Create: `.agent/verification.md`
- Modify: `AGENTS.md`
- Test: `scripts/check-architecture-docs.ts`

- [ ] **Step 1: Extend the documentation check with failing assertions**

Add the three new files to `files` and require their purpose-specific terms.
For example, require `internal/generated/fx`, `internal/lib`, and `samber/lo`
from the repository map; `OperationID`, `one command`, and `idempotent` from
the workflows; and `task sync`, `task quality`, and `task build` from the
verification matrix.

```ts
const guidanceFiles = [
  ".agent/repository-map.md",
  ".agent/workflows.md",
  ".agent/verification.md",
];
```

- [ ] **Step 2: Run the focused check and verify it fails**

Run: `rtk bun run check:architecture-docs`

Expected: non-zero exit with a missing guidance-file or missing-term error.

- [ ] **Step 3: Write the three focused references and update the index**

Write the files as reference/how-to documents, not a second README:

- `repository-map.md` lists ownership and source of truth for Ent schema,
  Goose migrations, Huma/OpenAPI, Hey API output, Fx registries, River,
  config, realtime, and web. It lists every generated surface with its
  Taskfile regeneration command. It defines `internal/platform` versus narrow
  tested `internal/lib/topic-name` packages and the direct-use policy for
  `github.com/samber/lo`.
- `workflows.md` has short steps for a write use case, HTTP operation,
  job/scheduler, Ent/Goose rollout, OpenAPI/web contract change, and local
  `task dev`. It states one mutating handler calls one command; a multi-write
  intent is one transaction-owning orchestrator, not calls to independently
  transactional `Execute` methods. It includes error-code, compatibility,
  idempotency, input-validation, ownership, and sensitive-data checks.
- `verification.md` maps each change category to focused tests plus `task
  sync`, `task quality`, and `task build` as required.
- `AGENTS.md` points agents to these documents after its existing instruction
  to read relevant `.agent` rules.

Include this exact type-safety stance in the new references:

```text
Public Huma DTOs are named structs with explicit fields and JSON tags.
Handwritten TypeScript uses generated Hey API types; boundary data is unknown
until narrowed. Do not introduce any, interface{}, or map[string]any by
default. Existing SystemMetadata JSON remains a documented schema-less
exception until product keys and values are specified.
```

- [ ] **Step 4: Run the focused check and verify it passes**

Run: `rtk bun run check:architecture-docs`

Expected: exit 0.

- [ ] **Step 5: Commit the documentation foundation**

```bash
rtk git add AGENTS.md .agent scripts/check-architecture-docs.ts
rtk git commit -m "docs: add agent development guidance"
```

## Task 2: Align focused rules with typed and safe delivery boundaries

**Files:**
- Modify: `.agent/architecture.md`
- Modify: `.agent/api.md`
- Modify: `.agent/database.md`
- Modify: `.agent/worker.md`
- Modify: `.agent/web.md`
- Modify: `.agent/config.md`
- Test: `scripts/check-architecture-docs.ts`

- [ ] **Step 1: Add failing rule-presence assertions**

Extend `requiredByFile` so each relevant focused rule contains its boundary:

```ts
".agent/api.md": ["OperationID", "one command", "strict Go structs"],
".agent/database.md": ["field.Enum", "date-only", "Schema-less JSON"],
".agent/worker.md": ["idempotent", "versioned", "retry"],
".agent/web.md": ["unknown", "any", "Hey API"],
```

Also reject `map[string]any` in the API guidance as an unqualified default;
do not reject it from the documented SystemMetadata exception.

- [ ] **Step 2: Run the check and verify it fails**

Run: `rtk bun run check:architecture-docs`

Expected: non-zero exit naming at least one missing new boundary term.

- [ ] **Step 3: Add narrow rules to their owners**

- API: named strict Huma DTOs, deterministic operation metadata, stable error
  codes, additive OpenAPI review, and exactly one mutating command.
- Database: typed Ent enum/JSON policies; UTC timestamp versus distinct
  date-only/time-only domain type; expand/migrate/contract rollout and
  explicit backfill.
- Worker: typed versioned arguments, idempotency, retry/timeout/cancellation,
  scheduler-enqueue-only ownership.
- Web: use generated Hey API types, narrow `unknown`, do not hand-edit
  generated clients or introduce handwritten `any`.
- Config: least-data boundary and secrets constraints.
- Architecture: dependency direction and prohibition on using command calls to
  bypass transactional orchestration.

Do not alter `SystemMetadata` or introduce authentication behaviour; those
need product requirements.

- [ ] **Step 4: Run the focused check and verify it passes**

Run: `rtk bun run check:architecture-docs`

Expected: exit 0.

- [ ] **Step 5: Commit the aligned focused rules**

```bash
rtk git add .agent scripts/check-architecture-docs.ts
rtk git commit -m "docs: define typed delivery guardrails"
```

## Task 3: Derive deterministic Huma domain tags

**Files:**
- Modify: `internal/transport/http/method/register.go`
- Modify: `internal/transport/http/method/register_test.go`

- [ ] **Step 1: Write failing tests for the tag helper and route metadata**

Add tests that require `/api/v1/system-metadata/{key}` to yield exactly
`[]string{"system-metadata"}` and that inspect the registered OpenAPI
operation tags:

```go
func TestTagsDerivesFirstV1Resource(t *testing.T) {
    if got := method.Tags("/api/v1/system-metadata/{key}"); !reflect.DeepEqual(got, []string{"system-metadata"}) {
        t.Fatalf("Tags() = %#v", got)
    }
}

func TestPUTSetsDomainTag(t *testing.T) {
    router := chi.NewRouter()
    api := humachi.New(router, huma.DefaultConfig("test", "1.0.0"))
    path := "/api/v1/examples/{id}"
    method.PUT(api, path, "Set example", method.Tags(path), func(context.Context, *putInput) (*putOutput, error) {
        return &putOutput{}, nil
    })
    operation := api.OpenAPI().Paths[path]
    if operation == nil || operation.Put == nil || !reflect.DeepEqual(operation.Put.Tags, []string{"examples"}) {
        t.Fatalf("operation = %#v", operation)
    }
}
```

Add a panic assertion for paths without a static resource after `/api/v1/` so
an invalid handwritten registration fails during application setup rather than
silently producing an ungrouped OpenAPI operation.

- [ ] **Step 2: Run the method package test and verify it fails**

Run: `rtk go test ./internal/transport/http/method -run 'Test(Tags|PUTSetsDomainTag)' -count=1`

Expected: FAIL because `method.Tags` does not exist.

- [ ] **Step 3: Implement the minimal tag helper**

Add an exported `Tags(path string) []string` in `register.go`. It must:

1. split the path into slash-delimited segments;
2. require the prefix `api/v1`;
3. return a single-element slice containing the first non-empty, non-parameter
   segment after that prefix;
4. panic with the path in the message if no resource segment exists.

Do not change the existing deterministic `OperationID` algorithm. Keep method
helpers accepting a tag slice, because handwritten operations may explicitly
provide the one derived tag.

- [ ] **Step 4: Run focused tests and verify they pass**

Run: `rtk go test ./internal/transport/http/method -count=1`

Expected: PASS, including existing typed-route and `OperationID` tests.

- [ ] **Step 5: Commit Huma tag derivation**

```bash
rtk git add internal/transport/http/method/register.go internal/transport/http/method/register_test.go
rtk git commit -m "feat: derive Huma resource tags"
```

## Task 4: Make HTTP scaffolds typed and tagged

**Files:**
- Modify: `vpkg/vandor/go/templates/http_handler.vxt`
- Modify: `vpkg/vandor/go/README.md`
- Create: `scripts/check-vpkg-templates.ts`
- Modify: `Taskfile.yml`
- Test: `rtk bun run scripts/check-vpkg-templates.ts` and `rtk vx view`

- [ ] **Step 1: Write a failing template-contract check**

Create `scripts/check-vpkg-templates.ts` that reads
`vpkg/vandor/go/templates/http_handler.vxt` and fails unless it contains the
named DTO declarations and the `method.Tags` call. Add it to the `test` task
after the existing documentation check. Its assertions are:

```ts
const source = await Bun.file("vpkg/vandor/go/templates/http_handler.vxt").text();
for (const fragment of [
  "type input struct{}",
  "type output struct{}",
  'method.Tags("{{ path }}")',
]) {
  if (!source.includes(fragment)) throw new Error(`HTTP handler template is missing ${fragment}`);
}
```

The check must fail against the current `nil` tag and anonymous `struct{}`
payloads. The test may validate template source; the next step validates vx
rendering without creating a source file.

- [ ] **Step 2: Run the focused check and verify it fails**

Run: `rtk bun run scripts/check-vpkg-templates.ts`

Expected: FAIL because the current generated source passes `nil` and uses
anonymous structs.

- [ ] **Step 3: Update the template and package documentation**

Generate private named input/output DTO declarations and use them in `Handle`.
Replace the `nil` argument with:

```go
method.Tags("{{ path }}")
```

Keep the public Taskfile arguments unchanged: only `name`, `method`, and
`path`. Update `vpkg/vandor/go/README.md` to say the path derives the one
kebab-case resource tag.


- [ ] **Step 4: Run template and sync checks**

Run:

```bash
rtk bun run scripts/check-vpkg-templates.ts
rtk vx view vandor/go:http-handler --plan --set name=SetSystemMetadata --set method=PUT --set path=/api/v1/system-metadata/{key}
rtk task sync
rtk go test ./vpkg/vandor/go/tools/sync ./internal/transport/http/method -count=1
```

Expected: all commands exit 0 and no generated Fx registry diff is introduced.

- [ ] **Step 5: Commit the scaffold change**

```bash
rtk git add Taskfile.yml scripts/check-vpkg-templates.ts \
  vpkg/vandor/go/templates/http_handler.vxt vpkg/vandor/go/README.md
rtk git commit -m "feat: scaffold tagged typed HTTP handlers"
```

## Task 5: Verify end-to-end guidance and preserve scope boundaries

**Files:**
- Verify only: documentation, Huma helper/template, Fx/OpenAPI/Hey generated
  output

- [ ] **Step 1: Confirm the SystemMetadata exception is not changed**

Run:

```bash
rtk git diff -- database/schema/systemmetadata.go internal/usecase/system_metadata.go \
  internal/transport/http/handler/system_metadata/set_system_metadata.go
```

Expected: no product-schema change unless the user separately supplies the
metadata key/value contract.

- [ ] **Step 2: Run full repository verification**

Run:

```bash
rtk task quality
rtk task build
rtk git diff --check
```

Expected: exit 0. Generated OpenAPI, Hey API, and Fx output are current.

- [ ] **Step 3: Inspect working tree before final commit**

Run: `rtk git status --short`

Expected: only files from this plan are staged or unstaged. Preserve the
pre-existing `apps/web/src/routeTree.gen.ts` change; do not stage, revert, or
claim it as part of this work.

- [ ] **Step 4: Report the checkpoint commits and remaining user change**

Report the commits from Tasks 1–4 and explicitly name the pre-existing
`apps/web/src/routeTree.gen.ts` modification. Do not stage, revert, or claim
that file as part of this work.

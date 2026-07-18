# Agent Development Guidance Design

**Status:** Proposed — awaiting implementation-plan review  
**Date:** 2026-07-19

## 1. Purpose

Give AI agents a compact, accurate, and executable understanding of the
repository. The guidance must reduce incorrect architectural changes without
duplicating the README, generated contracts, or implementation details that
belong in source code.

## 2. Documentation Shape

`AGENTS.md` remains the short entry point. It links agents to the focused rules
that match the work they are about to perform:

```text
.agent/
  repository-map.md  # reference: ownership and source of truth
  workflows.md       # how-to: common development changes
  verification.md    # reference: change-to-verification matrix
```

The existing architecture, API, database, worker, web, configuration, and UI
rules remain focused on their individual boundaries. The new files do not
become historical design archives.

## 3. Repository Map

The repository map identifies the source of truth, generated output, and
ownership for Ent schemas and migrations, Huma/OpenAPI and Hey API, Fx
registries, commands and adapters, River, configuration, realtime, and web.

Generated output is never hand-edited. The map states the exact Taskfile
command that regenerates each output. It distinguishes `internal/platform`
(infrastructure and external integration) from `internal/lib` (small,
project-owned reusable logic).

## 4. Development Workflows

The workflow guide covers use cases, HTTP operations, River jobs and schedules,
Ent schema/migration changes, API/OpenAPI/web changes, and local development.
Each workflow starts with the source of truth, gives the permitted dependency
direction, and finishes with its required sync and verification commands.

Use cases own write-side business rules and their Ent mutation plus River
enqueue transaction. A mutating Huma handler calls one command only. If one
HTTP operation needs several writes, it calls one intent-specific orchestrator
use case that owns one transaction; it must not call several independently
transactional command `Execute` methods. A handler may compose read-only data
when building a read model.

## 5. Huma Operation Metadata

Every business Huma operation has deterministic metadata:

- `OperationID` is generated from its HTTP method and full path.
- `Summary` is a short action phrase for humans.
- Exactly one tag represents the operation's domain/resource in kebab-case.

Tags are derived for generated handlers from the first static resource segment
after `/api/v1/`. For example,
`PUT /api/v1/system-metadata/{key}` uses `system-metadata`. The public
`task add:http-handler` interface remains limited to `name`, `method`, and
`path`; tags are not a fourth free-form scaffold input. Handwritten handlers
follow the same rule. The generator and its tests must prevent new scaffolds
from registering `nil` or empty tags.

## 6. Shared Helpers and External Utilities

Reusable, project-owned, context-free helpers live in narrow packages under
`internal/lib/<topic>`. A helper moves there only when it has a clear stable
abstraction and serves at least two contexts. Feature-local helpers remain
private and close to their feature.

`internal/lib` must not import use cases, handlers, jobs, schedulers, Ent,
River, Fx, or application configuration. It has tests for each package.
`internal/platform` remains the home for infrastructure integration.

`github.com/samber/lo` is an approved direct utility dependency for clear,
generic collection and optional-value transformations. Prefer the standard
library when equally clear. Do not create one-to-one wrappers around `lo`, use
panic helpers in long-lived runtimes, or introduce `lo/parallel` without an
explicit concurrency design. Add the dependency only at its first real use.

## 7. Type-Safe Contracts

Types describe data semantics before storage mechanics.

- Huma inputs and outputs use named, strict Go structs with explicit fields and
  JSON tags. Public HTTP contracts do not use `any`, `interface{}`, or
  `map[string]any` by default.
- TypeScript handwritten code does not use `any`. Boundary data begins as
  `unknown` and is narrowed before use. Generated Hey API output is not edited
  by hand; generator-originated `any` is handled by generator configuration or
  a narrowly scoped lint exception.
- Ent schemas model enums with `field.Enum` and named Go types when shared
  semantics warrant them. Timestamps use UTC `time.Time`; date-only and
  time-only values use distinct domain types rather than a timestamp with a
  dummy component.
- JSON columns use typed structs or typed slices. Schema-less JSON is an
  explicit exception, never a default convenience.
- HTTP DTOs are separate from generated Ent entities. API compatibility and
  storage representation may evolve independently.

The baseline `SystemMetadata.value` is currently schema-less JSON. Its
refactor to typed metadata must define the supported metadata keys and value
shapes before changing the Ent schema, use-case input, Huma input, OpenAPI, and
Hey API client together. This design does not invent those product schemas.

## 8. Error Contracts and API Compatibility

Huma operations use one documented error-envelope convention with stable,
machine-readable codes. Handlers map domain errors to that convention rather
than inventing transport-specific response shapes. Error details do not expose
secrets, credentials, SQL statements, or internal topology.

Public OpenAPI changes are additive by default. Removing or renaming a public
field, changing its meaning or type, tightening a previously accepted payload,
or changing a success/error status is a compatibility review. The workflow
identifies the Huma DTO, OpenAPI document, generated Hey API client, and web
consumers that must change together.

## 9. Ent, Goose, and Data Rollout

Schema changes follow a safe order: update the Ent schema, generate the Ent
client, add an explicit Goose migration, test the migration path, and keep the
application compatible with the deployed data during rollout. Data backfills
are explicit, retryable operations rather than incidental side effects of an
application startup.

Enum, temporal, and typed-JSON changes receive a compatibility review. The
workflow states when an expand/migrate/contract rollout is required instead of
assuming a destructive schema change is safe.

## 10. River Reliability and Async Contracts

River job arguments are typed, serializable contracts. They are versioned when
their meaning changes while queued work may still exist. Jobs and periodic
schedules are idempotent: retries, duplicate execution, and horizontally
replicated scheduler processes must not produce duplicate business effects.

Each job explicitly chooses retry, timeout, cancellation, and failure-handling
behavior. A schedule only enqueues work; it never becomes the source of truth
for deadline-sensitive business state. Mutation use cases remain responsible
for deciding which River work is enqueued inside their transaction.

## 11. Security and Boundary Data

Every new public operation identifies its input validation, authorization or
ownership boundary, and sensitive fields before implementation. Browser,
realtime, job, log, error, and configuration boundaries receive only the data
they need. Private credentials, database connection details, and internal
realtime credentials never cross into browser-visible configuration or public
responses.

The guidance does not define a product authentication model. It prevents an
agent from bypassing or inventing one while adding a feature.

## 12. Verification

The verification matrix maps change categories to focused tests, `task sync`,
`task quality`, and `task build`. Existing architecture-document checks expand
to require the stable guidance entry points and prevent regressions to the
documented Huma, type-safety, command-boundary, and generator rules.

## 13. Non-Goals

- Do not create a catch-all `utils` package.
- Do not turn use-case composition into a workaround for transaction design.
- Do not add a free-form HTTP tag argument to Taskfile scaffolding.
- Do not force a product schema for existing system metadata without product
  requirements.
- Do not modify generated output by hand.

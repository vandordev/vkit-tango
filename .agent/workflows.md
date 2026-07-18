# Development Workflows

## Write use case

1. Start with a focused test in `internal/usecase/<intent>_test.go`.
2. Add one public command struct in `internal/usecase/<intent>.go`; writes own
   one Ent/SQL transaction and any River enqueue inside it.
3. Run the focused Go test, then `task sync:usecase`.
4. A multi-write intent is one transaction-owning orchestrator. Do not call
   independently transactional command `Execute` methods from each other.

## Huma HTTP operation

1. One method plus path is one handler file. A mutating handler validates a
   named, strict Go struct then calls one command.
2. Public Huma DTOs have explicit fields and JSON tags. Do not introduce
   `any`, `interface{}`, or `map[string]any` by default. The existing
   SystemMetadata JSON value is a documented schema-less exception until its
   product keys and value shapes are defined.
3. Use the method helper with a deterministic `OperationID`, a concise action
   summary, and exactly one kebab-case resource tag. The tag is derived from
   the first resource segment after `/api/v1/`.
4. Map domain errors to the shared stable error-code envelope. Do not expose
   secrets, SQL, credentials, or internal topology. Treat non-additive
   OpenAPI changes as a compatibility review.
5. Run `task sync:http`, then regenerate the OpenAPI and Hey API client with
   `task sync` when the public contract changes.

## River job and scheduler

1. Jobs decode typed, serializable, versioned arguments and invoke commands.
2. Choose idempotent behaviour, retry limits, timeout, cancellation, and
   failure handling explicitly. Retries and duplicate delivery must not repeat
   business effects.
3. Schedulers only enqueue jobs; they never mutate domain state directly.
4. Run `task sync:worker` or `task sync:scheduler` after adding the adapter.

## Ent, Goose, and data rollout

1. Model semantics before storage: `field.Enum` for closed values, UTC
   `time.Time` for timestamps, and distinct domain types for date-only or
   time-only values. JSON columns use typed structs or slices.
2. Change Ent schema, run `task db:generate`, add an explicit Goose migration,
   and test the migration path. For destructive or typed-data changes, use an
   expand/migrate/contract rollout with an explicit retryable backfill.
3. Do not copy generated Ent entities directly into public HTTP DTOs.

## Web, realtime, and boundary data

1. Use Hey API generated types in handwritten TypeScript. Boundary data starts
   as `unknown` and is narrowed before use; handwritten `any` is not allowed.
2. Do not edit generated clients. Solve generator-originated `any` through a
   generator change or a narrowly scoped lint exception.
3. The shared QueryClient invalidates all query keys after a successful
   mutation without awaiting refetches. Do not repeat global invalidation in
   each mutation; use a local awaited invalidation only when a specific UX
   must wait for fresh data. Socket.IO events still invalidate affected query
   keys for mutations made outside this browser.
4. Validate public input and identify authorization or ownership before adding
   an operation. Browser, realtime, job, log, error, and config boundaries
   receive only needed data. Private credentials never reach browser config or
   public responses.

## Local development

Run `task migrate` before long-lived services. `task dev` starts API, worker,
scheduler, and web; select services with `task dev -- api web realtime`.
Realtime is opt-in.

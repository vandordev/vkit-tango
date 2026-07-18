# Agent Development Pain Points

This is the evidence-backed feedback loop for improving the baseline. It is a
maintenance backlog, not a project instruction file and not a historical
design archive. Consult it when a task exposes friction in shared tooling,
generators, runtime wiring, or agent guidance.

## When to add an entry

Add a pain point only after it occurs in this repository. Do not record a
hypothetical concern. First search existing entries and update the same issue
instead of duplicating it.

Every entry must include reproducible evidence, the user or agent impact, a
temporary workaround, one concrete baseline improvement, and a Regression
check. Avoid entries that merely say an agent was confused.

## Entry template

```md
### AP-### — Short problem statement

- Category: Generator | Runtime | Tooling | Documentation | Contract | Testing
- Context: Feature, command, or development path where it occurred.
- Evidence: Exact command, error, file path, or reproducible observation.
- Impact: What work becomes slow, unsafe, repetitive, or ambiguous.
- Workaround: Safe temporary path until the baseline changes.
- Proposed baseline improvement: Smallest durable change that removes the friction.
- Regression check: Exact automated test or command that prevents recurrence.
- Status: Open | Accepted | Resolved | Declined
- Resolution: Commit, decision, or reason; required when status is Resolved or Declined.
```

## Open

No open pain points are currently recorded.

## Resolved

### AP-001 — Generated TanStack Start route tree was tracked

- Category: Tooling
- Context: Web builds regenerated `apps/web/src/routeTree.gen.ts`, creating
  unrelated working-tree diffs during backend and documentation work.
- Evidence: `task build` included a Git diff guard for this generated file.
- Impact: Generated route changes could be mistaken for user-authored changes.
- Workaround: Regenerate before inspecting the working tree.
- Proposed baseline improvement: Ignore and untrack the generated route tree.
- Regression check: `git check-ignore -v --no-index apps/web/src/routeTree.gen.ts`.
- Status: Resolved
- Resolution: `6f6d975 chore: ignore generated route tree`.

### AP-002 — HTTP scaffold omitted domain metadata and named DTOs

- Category: Generator
- Context: `task add:http-handler` produced a handler with a nil Huma tag and
  anonymous input/output structs.
- Evidence: `vpkg/vandor/go/templates/http_handler.vxt` passed `nil` to the
  method helper and used `struct{}` in `Handle`.
- Impact: Generated OpenAPI operations lacked consistent resource grouping and
  scaffolded contracts were less explicit.
- Workaround: Edit the new handler manually after generation.
- Proposed baseline improvement: Derive a resource tag from the `/api/v1` path
  and generate named private DTOs.
- Regression check: `rtk bun run scripts/check-vpkg-templates.ts` and
  `rtk go test ./internal/transport/http/method -count=1`.
- Status: Resolved
- Resolution: `25398a2 feat: scaffold tagged typed HTTP handlers`.

## Review lifecycle

Review Open entries before a baseline release or after completing a feature
that touched the reported area. A Resolved entry stays as a short learning
record with its commit and Regression check. Mark an entry Declined only with
an explicit reason so future agents do not reopen the same discussion.

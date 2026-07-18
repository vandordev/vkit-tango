# AGENTS.md

Read `README.md`, the relevant `.agent` rules, `.agent/repository-map.md`,
`.agent/workflows.md`, `.agent/verification.md`, and `git status --short`
before editing.

The web uses shadcn/ui as its one primary UI system; Mantine and MUI are
alternatives that must be selected explicitly.

- `apps/api`, `apps/worker`, and `apps/scheduler` are thin Uber Fx composition roots. `apps/migrate` remains a short-lived non-Fx migration command.
- Go owns Huma HTTP, Ent persistence, Goose migrations, River jobs/schedules, and every business mutation.
- Active Huma routes use `/api/v1/*`; health/docs/OpenAPI are unversioned process routes.
- Reads may use Ent in their owning HTTP handler. Writes must call `internal/usecase`; a use case owns its Ent/SQL transaction and any River enqueue.
- Adapters depend on shared `internal/contract` boundaries: commands for mutations, HTTP handlers for routes, River jobs for worker registration, and schedulers for periodic registration. Jobs execute commands; schedules only enqueue jobs.
- Generate scaffolds and registries only through Taskfile: `task add:*` and `task sync*`. Do not invoke `vx` directly or hand-edit `internal/generated/fx`.
- Ent schema lives in `database/schema`; generated Ent output is `internal/platform/db`; migrations use Goose.
- TypeScript owns only TanStack Start and Socket.IO. Hey API generates the web client from `contracts/openapi/openapi.json`; realtime is documented in `contracts/asyncapi/realtime.v1.yaml`.
- Use the shared YAML modules under `config`; never expose database or private realtime credentials to the browser.
- Run focused tests first. Run `task sync` after generator or registry changes; run `task quality` and `task build` for shared/runtime changes.

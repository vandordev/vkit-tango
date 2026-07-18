# AGENTS.md

Read `README.md`, the relevant `.agent` rules, and `git status --short` before editing.

Web memakai shadcn/ui sebagai one primary UI system; Mantine dan MUI adalah alternatif yang harus dipilih secara eksplisit.

- Go owns Huma HTTP, Ent persistence, Goose migrations, River schedules/jobs, and every business mutation.
- Active Huma routes use `/api/v1/*`; health/docs/OpenAPI are unversioned process routes.
- Reads may use Ent in their owning HTTP handler. Writes must call `internal/usecase`; a usecase owns its Ent/SQL transaction and any River enqueue.
- Ent schema lives in `database/schema`; generated Ent output is `internal/platform/db`; migrations use Goose.
- TypeScript owns only TanStack Start and Socket.IO. Hey API generates the web client from `contracts/openapi/openapi.json`; realtime is documented in `contracts/asyncapi/realtime.v1.yaml`.
- Use the shared YAML modules under `config`; never expose database or private realtime credentials to the browser.
- Run focused tests first, then `task quality` and `task build` for shared/runtime changes.

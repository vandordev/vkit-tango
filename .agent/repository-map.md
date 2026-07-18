# Repository Map

Use this reference to find the source of truth before editing. Do not hand-edit
generated output.

| Concern | Source of truth | Generated/output | Refresh command |
| --- | --- | --- | --- |
| Ent persistence | `database/schema` | `internal/platform/db` | `task db:generate` |
| Schema history | `database/migrations` | PostgreSQL schema | `task migrate` |
| HTTP contract | Huma handlers under `internal/transport/http/handler` | `contracts/openapi/openapi.json` | `task api:openapi` |
| Web API client | Huma OpenAPI | `apps/web/src/lib/api/generated` | `task api:client:generate` |
| Fx registration | commands and adapter contract implementations | `internal/generated/fx` | `task sync` or targeted `task sync:*` |
| HTTP scaffold | `vpkg/vandor/go/templates/http_handler.vxt` | new handler source | `task add:http-handler` |
| River scaffold | `vpkg/vandor/go/templates/{job,scheduler}.vxt` | new job/scheduler source | `task add:job` or `task add:scheduler` |
| Realtime contract | `contracts/asyncapi/realtime.v1.yaml` | TypeScript Socket.IO runtime | project-specific tests |
| Runtime config | shared YAML under `config` | typed Go/TypeScript config | no generated source |

`apps/api`, `apps/worker`, and `apps/scheduler` are distinct Fx composition
roots. `apps/migrate` remains non-Fx. `internal/contract` is the shared
boundary for commands, HTTP handlers, jobs, and scheduler registrars.

`internal/platform` owns infrastructure and external integration: database,
River, realtime publishing, and runtime wiring. `internal/lib/topic-name` is
only for narrow, context-free project helpers that serve at least two contexts.
Every `internal/lib` package has tests and must not import use cases, handlers,
jobs, schedulers, Ent, River, Fx, or application configuration. Keep
feature-local helpers private beside the feature.

`github.com/samber/lo` is approved for clear generic collection and optional
transformations. Prefer the standard library when equally clear, use `lo`
directly instead of one-to-one wrappers, and add it only at its first real use.
Do not use panic helpers in long-lived runtimes or `lo/parallel` without an
explicit concurrency design.

Application developers and CI use Taskfile. `vx` is an internal detail of
`task add:*`; maintainers may preview templates with `vx view ... --plan`.

# Verification Matrix

Run the focused test first. The matrix adds the smallest required shared checks.

| Change | Focused verification | Required follow-up |
| --- | --- | --- |
| Use case | package Go test | `task sync:usecase`; `task quality` when shared/runtime code changes |
| HTTP handler/OpenAPI | handler or method Go test | `task sync`; web type-check when client consumers change |
| Job or scheduler | adapter/package Go test | targeted `task sync:*`; `task quality` for runtime changes |
| Fx registry or generator | generator/template test | `task sync`, `task quality`, `task build` |
| Ent schema or Goose migration | schema/migration test | `task db:generate`, migration check, `task quality`, `task build` |
| Web/realtime TypeScript | focused Bun test | lint and type-check; `task quality` for shared contract changes |
| Config, runtime, or composition root | focused Go test | `task quality` and `task build` |
| Documentation or Taskfile surface | check script | `task test`; `task quality` when the command surface changes |

`task sync` refreshes Fx registries, OpenAPI, and the Hey API client. Generated
output must be committed when it changes. Before claiming completion, run
`rtk git diff --check` and inspect `rtk git status --short` for unrelated user
changes.

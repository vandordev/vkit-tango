# Architecture Rules

`apps/api`, `apps/worker`, and `apps/scheduler` are distinct Uber Fx composition roots; `apps/migrate` is the short-lived non-Fx Goose/River migration command. Huma is the only HTTP server and active business routes use `/api/v1`. Ent is the only ORM, Goose is the application migration mechanism, and River is the PostgreSQL queue.

Go use cases own every write and their Ent mutation plus River enqueue share one SQL transaction; direct Ent reads are allowed only for read models. Shared `internal/contract` interfaces connect commands, HTTP handlers, River jobs, and periodic schedulers. Jobs execute commands, while scheduler registrations only enqueue typed jobs. Generated `internal/generated/fx` registries wire documented contract implementations without runtime reflection. Hey API generates the web contract from Huma OpenAPI. Socket.IO stays in TypeScript and receives authenticated, versioned realtime events from Go.

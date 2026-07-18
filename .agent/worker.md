# Worker Rules

`apps/worker` is the Uber Fx root that executes Go-only River jobs. A job decodes typed, versioned arguments and invokes its matching `internal/contract.Command`; it may use Ent for reads but never restates mutation rules. Job behaviour is idempotent across retries and duplicate delivery, with retry, timeout, cancellation, and failure handling chosen explicitly. Worker registrations are generated from `internal/contract.Job` implementations and refreshed with `task sync:worker` or `task sync`.

`apps/scheduler` is a separate Uber Fx process that owns periodic River enqueue registrations. A schedule only enqueues typed jobs and never invokes a use case directly; horizontal safety relies on River's supported periodic scheduling behavior. Huma owns `/api/v1`; Goose controls schema; Hey API and Socket.IO remain external transport boundaries. Realtime delivery uses authenticated HTTP and at-least-once retries.

# Worker Rules

River workers are Go-only and register periodic reconciliation jobs on every replica. A job may use Ent for reads but invokes a Go usecase for writes. Huma owns `/api/v1`; Goose controls schema; Hey API and Socket.IO remain external transport boundaries. Realtime delivery uses authenticated HTTP and at-least-once retries.

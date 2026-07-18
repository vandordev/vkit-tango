# Shared runtime configuration

Each runtime selects only the YAML modules it needs. Keys are snake_case and
secrets are interpolation references, never literal values. Supported
interpolation is `${NAME}` and `${NAME:-fallback}` only.

Go API selects `app`, `database`, `http_api`, `realtime`, and `observability`.
Go worker selects `app`, `database`, `worker`, `realtime`, and `observability`.
The migration runtime selects `database`. The TypeScript web runtime receives
only `web.public`; database and realtime internal secrets stay server-only.

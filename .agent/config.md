# Configuration Rules

Shared YAML uses snake_case and `${NAME}` / `${NAME:-fallback}` only. Huma, Ent, Goose, and River load their Go subsets independently from Hey API/web configuration. `/api/v1` remains browser-safe through the public web subtree. Socket.IO internal credentials never enter browser configuration.

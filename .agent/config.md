# Configuration Rules

Shared YAML uses snake_case and `${NAME}` / `${NAME:-fallback}` only. Fx modules load typed Go subsets for Huma, Ent, Goose, River worker, and scheduler infrastructure independently from Hey API/web configuration. Apply least data at every config boundary: `/api/v1` remains browser-safe through the public web subtree, while database credentials, private River settings, and Socket.IO internal credentials never enter browser configuration.

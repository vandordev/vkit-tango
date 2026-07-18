# API Rules

Huma Go handlers own `/api/v1/*`. Query handlers may read Ent directly; mutation handlers validate transport input and invoke an `internal/usecase` command. Keep process health, docs, and OpenAPI outside the versioned business prefix.

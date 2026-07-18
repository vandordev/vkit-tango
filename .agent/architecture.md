# Architecture Rules

Huma is the only HTTP server; active business routes use `/api/v1`. Ent is the only ORM and Goose is the only application migration mechanism. River is the PostgreSQL queue and its schedules run inside the Go worker. Go usecases own every write, while direct Ent reads are allowed only for read models. Hey API generates the web contract from Huma OpenAPI. Socket.IO stays in TypeScript and receives authenticated, versioned realtime events from Go.

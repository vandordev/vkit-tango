# Database Rules

Ent schema source is `database/schema`, generated output is `internal/platform/db`, and Goose migrations are the deployed schema history. Go usecases own Ent writes and River enqueueing in one SQL transaction. Huma exposes writes under `/api/v1`; Hey API transports reads to web; Socket.IO only signals refetches.

# API Rules

Huma handlers register business operations only through `/api/v1`. A mutation validates transport input then invokes an intent-named Go usecase; an Ent read may shape a response directly. Goose owns schema history, River owns asynchronous work, and Hey API consumes this Huma contract. Socket.IO is only a realtime invalidation boundary and never a write path.

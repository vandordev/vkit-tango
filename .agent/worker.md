# Worker Rules

River workers are Go-only. They execute typed jobs, retry failed delivery, and invoke usecases where applicable. Product rules belong in usecases, not job handlers. Mutations that enqueue work must share one PostgreSQL transaction.

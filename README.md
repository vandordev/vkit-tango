# vkit-fast

Monorepo data-driven untuk TanStack Start, Go, PostgreSQL, dan Socket.IO.

UI web memakai shadcn/ui sebagai one primary UI system; Mantine atau MUI hanya dipilih secara sengaja per proyek.

## Architecture

- `apps/web`: TanStack Start. Browser memanggil Go API melalui `/api/*` dan menggunakan client serta TanStack Query hooks hasil Hey API.
- `apps/api`: HTTP API Go (Huma). Semua endpoint bisnis menggunakan `/api/v1/*`; `/health`, `/health/ready`, `/api/openapi.json`, dan `/api/docs` bersifat process-level.
- `apps/worker`: River worker Go. Semua background processing dan jadwal River berjalan di Go.
- `apps/migrate`: satu proses Goose dan River migration.
- `database/schema`: Ent schema; generated client berada di `internal/platform/db`.
- `database/migrations`: Goose migrations.
- `internal/usecase`: satu sumber aturan mutasi bisnis. Handler query boleh memakai Ent langsung; mutasi wajib lewat usecase.
- `apps/realtime` dan `packages/realtime`: Socket.IO TypeScript. Go mempublikasikan event melalui endpoint privat `/internal/events`; kontraknya ada di `contracts/asyncapi/realtime.v1.yaml`.

PostgreSQL adalah satu-satunya queue backend. Mutasi Ent dan enqueue River harus berbagi transaksi SQL yang sama.

## Configuration

Konfigurasi YAML snake_case berada pada `config/`. Go memuat modul yang dibutuhkan secara eksplisit (`database`, `http_api`, `worker`, `realtime`), sedangkan TypeScript hanya memuat `web` atau `realtime`. Rahasia hanya melalui interpolasi environment.

## Commands

```bash
task install
cp .env.example .env
task migrate
task dev:api
task dev:worker
task dev:web
task dev:realtime
task api:client:generate
task quality
task build
```

Use `docker compose up --build` to run PostgreSQL, migrations, and the Go API. Add `--profile jobs` for the worker.

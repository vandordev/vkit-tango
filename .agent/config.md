# Configuration Rules

- Committed `config/*.yaml` modules are the primary configuration source. A runtime explicitly selects its modules and their merge order; module names do not imply runtime ownership.
- Merge selected modules left-to-right. Plain objects deep-merge; later arrays, scalars, and `null` replace earlier values.
- YAML may contain only `${NAME}` or `${NAME:-fallback}` interpolation. It resolves once against the process environment after merging. A missing or empty required value fails before Zod validation.
- Never place a secret literal in YAML. Developers copy `.env.example` to the one ignored root `.env`; production injects the same values through its deployment platform or secret manager.
- Server runtimes use the YAML wrapper and typed factories from `@repo/config` backed by `@t3-oss/env-core`. Add a key to the smallest runtime schema that needs it.
- TanStack Start runs through Vite and Nitro. The wrapper must run before `vite` development and production builds; browser code may receive only resolved `VITE_*` values.
- Embedded Elysia routes in TanStack Start require `DATABASE_URL` as a server-only web value. Never place it in browser code or expose it with `VITE_`.
- Optional S3/MinIO and realtime credentials stay server-only. Validate them through `createStorageConfig` and `createRealtimeConfig` respectively.

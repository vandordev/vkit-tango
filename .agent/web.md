# Web Rules

- Organize UI routes under `apps/web/src/routes`; `__root.tsx` owns the HTML shell and shared providers.
- Use Eden from `apps/web/src/lib/api`; the server client may use `treaty(app).api`, while browser calls use same-origin `/api`. Do not use tRPC, rewrites, or ad-hoc `fetch` calls in pages/components.
- Keep `src/server.ts` as the thin TanStack Start server entry: it forwards `/api/*` and `/health` to Elysia `app.fetch`, then delegates all other requests to Start.
- Browser-visible configuration uses an explicit `VITE_*` key. Database credentials and API secrets never belong in browser code.
- Keep API data loading in TanStack Router loaders, server functions, or typed client hooks; UI components do not import infrastructure packages.

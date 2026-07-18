# Web Rules

- Organize UI routes under `apps/web/src/routes`; `__root.tsx` owns the HTML shell and shared providers.
- Use the Hey API-generated client and TanStack Query hooks from the Huma OpenAPI contract. Browser requests use the same-origin `/api/*` boundary; do not use tRPC or ad-hoc API clients in pages/components.
- Keep the TanStack Start server entry thin and delegate Go API traffic to the configured upstream. The API's active business contract is `/api/v1`; health and OpenAPI/docs are process routes.
- Browser-visible configuration uses an explicit `VITE_*` key. Database credentials, River configuration, and Socket.IO private credentials never belong in browser code.
- Keep API data loading in TanStack Router loaders, server functions, or typed client hooks; UI components do not import infrastructure packages. Handwritten TypeScript uses Hey API types, treats boundary values as unknown until narrowed, and does not introduce `any`; generated client code is not edited manually. Socket.IO remains the realtime invalidation boundary, while Go/Ent/Goose own writes and River owns background work.

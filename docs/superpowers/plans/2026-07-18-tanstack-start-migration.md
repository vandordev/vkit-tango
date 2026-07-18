# TanStack Start Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the former `apps/web` runtime with a Bun-targeted TanStack Start application while retaining the embedded Elysia API and existing public endpoints.

**Architecture:** TanStack Start owns the application document, file-based page routes, and a thin custom server entry. `src/server.ts` forwards matching `/api/*` and `/health` Web `Request` objects to Elysia's existing `app.fetch`, while Eden continues to be the sole browser and server API client. Vite builds the web application and Nitro emits a Bun-compatible `.output` server.

**Tech Stack:** Bun 1.3.14, Vite 8.1.5, TanStack Start 1.168.30, TanStack Router 1.170.18, TanStack Router plugin 1.168.22, Nitro 3.0.260610-beta, React 19, shadcn/ui, Tailwind CSS v4, Elysia, Eden, Prisma, Turbo, Docker.

> **Superseded UI baseline:** The Mantine-specific dependency and root-route examples below record the migration’s original state. The current default is shadcn/ui; Mantine, MUI, and other libraries are valid alternatives only when a project deliberately selects one primary UI system. See [the shadcn/ui baseline design](../specs/2026-07-18-shadcn-ui-baseline-design.md).

---

## File structure

| Path | Responsibility |
| --- | --- |
| `apps/web/vite.config.ts` | Compose Vite, TanStack Start, Router generation, React, and Bun-targeted Nitro. |
| `apps/web/src/router.tsx` | Create and register the typed TanStack Router instance. |
| `apps/web/src/routes/__root.tsx` | HTML document, metadata, global shadcn/Tailwind styles, React Query, and fallback states. |
| `apps/web/src/routes/index.tsx` | Public `/` route. |
| `apps/web/src/routes/dashboard.tsx` | `/dashboard` route. |
| `apps/web/src/server.ts` | TanStack Start server entry that forwards Elysia `/api/*` and `/health` requests. |
| `apps/web/src/server/elysia-adapter.ts` | Server-only forwarding functions used by the route adapters and direct unit tests. |
| `apps/web/src/lib/api/client.ts` | Same-origin browser Eden facade. |
| `apps/web/src/lib/api/server.ts` | Server-only in-process Eden facade. |
| `apps/web/src/components/query-provider.tsx` | Per-renderer React Query provider retained from the existing web app. |
| `apps/web/src/styles.css` | Existing global styles moved from the Next app directory. |
| `Dockerfile.web` | Builds and runs the generated Bun/Nitro output. |

The migration deletes `apps/web/app/`, `apps/web/next.config.mjs`, `apps/web/next.config.test.ts`, `apps/web/postcss.config.mjs`, and `apps/web/lib/env.ts`. No generated `src/routeTree.gen.ts` is edited by hand.

### Task 1: Establish the Vite/TanStack Start build boundary

**Files:**
- Create: `apps/web/vite.config.ts`
- Create: `apps/web/vite.config.test.ts`
- Modify: `apps/web/package.json`
- Modify: `apps/web/tsconfig.json`
- Modify: `apps/web/eslint.config.js`
- Delete: `apps/web/next.config.mjs`
- Delete: `apps/web/next.config.test.ts`
- Delete: `apps/web/postcss.config.mjs`
- Modify: `turbo.json`
- Modify: `.gitignore`
- Modify: `bun.lock`

- [ ] **Step 1: Replace the old Next boundary test with a failing Vite boundary test.**

Create `apps/web/vite.config.test.ts`:

```ts
import { expect, test } from "bun:test";

test("configures TanStack Start with Bun-targeted Nitro", async () => {
  const viteConfig = await Bun.file(new URL("./vite.config.ts", import.meta.url)).text();

  expect(viteConfig).toContain('from "@tanstack/react-start/plugin/vite"');
  expect(viteConfig).toContain('from "@tanstack/router-plugin/vite"');
  expect(viteConfig).toContain('from "nitro/vite"');
  expect(viteConfig).toContain('nitro({ preset: "bun" })');
  expect(viteConfig).toContain('tanstackStart()');
  expect(viteConfig).toContain('tanstackRouter({ target: "react", autoCodeSplitting: true })');
});

test("uses the YAML wrapper without a Next.js command", async () => {
  const { scripts } = await Bun.file(new URL("./package.json", import.meta.url)).json();

  expect(scripts.dev).toContain("--env-file=../../.env");
  expect(scripts.dev).toContain("--modules base,web,api,storage");
  expect(scripts.dev).toContain("vite --port 4100");
  expect(scripts.build).toContain("vite build");
  expect(scripts.start).toContain(".output/server/index.mjs");
  expect(scripts.dev).not.toContain("next");
  expect(scripts.build).not.toContain("next");
});
```

- [ ] **Step 2: Run the test to verify it fails because the Vite configuration does not exist.**

Run: `rtk bun test apps/web/vite.config.test.ts`

Expected: FAIL with an error opening `apps/web/vite.config.ts` and with missing TanStack Start script assertions.

- [ ] **Step 3: Replace the web manifest, TypeScript setup, and lint setup.**

Set the relevant `apps/web/package.json` fields to:

```json
{
  "scripts": {
    "dev": "bun --env-file=../../.env run ../../packages/config/src/run.ts --modules base,web,api,storage -- bun run vite --port 4100",
    "build": "NODE_ENV=production bun --env-file=../../.env run ../../packages/config/src/run.ts --modules base,web,api,storage -- bun run vite build",
    "start": "bun --env-file=../../.env run ../../packages/config/src/run.ts --modules base,web,api,storage -- bun .output/server/index.mjs",
    "lint": "eslint",
    "check-types": "tsc --noEmit",
    "clean": "rm -rf .output .vite src/routeTree.gen.ts"
  },
  "dependencies": {
    "@elysia/eden": "^1.4.10",
    "@mantine/core": "^8.3.0",
    "@mantine/hooks": "^8.3.0",
    "@mantine/notifications": "^8.3.0",
    "@repo/api": "*",
    "@tanstack/react-query": "^5.101.2",
    "@tanstack/react-router": "1.170.18",
    "@tanstack/react-start": "1.168.30",
    "lucide-react": "^0.562.0",
    "nitro": "npm:nitro-nightly@3.0.1-20260717-080150-bfc2f5ef",
    "react": "^19.2.0",
    "react-dom": "^19.2.0",
    "zod": "^4.4.3"
  },
  "devDependencies": {
    "@repo/eslint-config": "*",
    "@repo/typescript-config": "*",
    "@tanstack/router-plugin": "1.168.22",
    "@types/node": "^22.15.3",
    "@types/react": "19.2.2",
    "@types/react-dom": "19.2.2",
    "@vitejs/plugin-react": "6.0.3",
    "elysia": "^1.4.21",
    "eslint": "^9.39.1",
    "typescript": "5.9.2",
    "vite": "8.1.5"
  }
}
```

Create `apps/web/vite.config.ts`:

```ts
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import { nitro } from "nitro/vite";

export default defineConfig({
  plugins: [
    tanstackRouter({ target: "react", autoCodeSplitting: true }),
    tanstackStart(),
    nitro({ preset: "bun" }),
    react(),
  ],
  resolve: {
    alias: {
      "@": new URL("./src", import.meta.url).pathname,
    },
  },
  ssr: {
    noExternal: ["@repo/api", "@repo/config", "@repo/database", "@repo/storage"],
  },
});
```

Replace `apps/web/tsconfig.json` with:

```json
{
  "extends": "@repo/typescript-config/base.json",
  "compilerOptions": {
    "baseUrl": ".",
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "jsx": "react-jsx",
    "noEmit": true,
    "paths": {
      "@/*": ["./src/*"]
    },
    "types": ["vite/client"]
  },
  "include": ["src/**/*.ts", "src/**/*.tsx", "vite.config.ts", "types/**/*.d.ts"],
  "exclude": ["node_modules", ".output"]
}
```

Replace `apps/web/eslint.config.js` with:

```js
import { config } from "@repo/eslint-config/base";

export default config;
```

Delete the three listed Next-specific configuration files. Update `turbo.json` so `globalEnv` contains only `NODE_ENV`, `LOG_LEVEL`, and `DATABASE_URL`; replace `.next/**` build output with `.output/**` and remove `NEXT_PUBLIC_APP_URL` from task environment entries. Add `.output/` and `.vite/` to `.gitignore`.

- [ ] **Step 4: Regenerate the lockfile without lifecycle scripts.**

Run: `rtk bun install --ignore-scripts`

Expected: PASS; `bun.lock` removes Next packages from the web dependency graph and records the pinned Vite, TanStack, and Nitro packages.

- [ ] **Step 5: Run the focused boundary test.**

Run: `rtk bun test apps/web/vite.config.test.ts`

Expected: PASS. The web typecheck runs after Task 3 creates the router and root route.

- [ ] **Step 6: Commit the build-boundary change.**

```bash
rtk git add apps/web/package.json apps/web/tsconfig.json apps/web/eslint.config.js apps/web/vite.config.ts apps/web/vite.config.test.ts turbo.json .gitignore bun.lock
rtk git rm apps/web/next.config.mjs apps/web/next.config.test.ts apps/web/postcss.config.mjs
rtk git commit -m "refactor(web): replace Next build boundary with Vite"
```

### Task 2: Add testable embedded Elysia adapters

**Files:**
- Create: `apps/web/src/server/elysia-adapter.ts`
- Create: `apps/web/src/server/elysia-adapter.test.ts`
- Create: `apps/web/src/server.ts`
- Create: `apps/web/src/server.test.ts`

- [ ] **Step 1: Write failing forwarding tests.**

Create `apps/web/src/server/elysia-adapter.test.ts`:

```ts
import { expect, test } from "bun:test";
import { resolvedConfigEnvironment } from "../../../../packages/config/src/run";

Object.assign(
  process.env,
  resolvedConfigEnvironment(["base", "api", "web"], {
    DATABASE_URL: "postgresql://db",
    NODE_ENV: "test",
  }),
);

test("forwards an API request to embedded Elysia", async () => {
  const { forwardApiRequest } = await import("./elysia-adapter");
  const response = await forwardApiRequest(new Request("http://localhost:4100/api/status"));

  expect(response.status).toBe(200);
  expect(await response.json()).toEqual({ success: true, data: { status: "ok" } });
  expect(response.headers.get("x-request-id")).toBeString();
});

test("forwards health requests to embedded Elysia", async () => {
  const { forwardHealthRequest } = await import("./elysia-adapter");
  const response = await forwardHealthRequest(new Request("http://localhost:4100/health"));

  expect(response.status).toBe(200);
  expect((await response.json()).data.status).toBe("healthy");
});
```

- [ ] **Step 2: Run the test to verify it fails because the adapter does not exist.**

Run: `rtk bun test apps/web/src/server/elysia-adapter.test.ts`

Expected: FAIL with `Cannot find module './elysia-adapter'`.

- [ ] **Step 3: Implement the adapters and server entry.**

Create `apps/web/src/server/elysia-adapter.ts`:

```ts
import { app } from "@repo/api";

export function forwardApiRequest(request: Request): Response | Promise<Response> {
  return app.fetch(request);
}

export function forwardHealthRequest(request: Request): Response | Promise<Response> {
  return app.fetch(request);
}
```

Create `apps/web/src/server.ts`:

```ts
import handler, { createServerEntry } from "@tanstack/react-start/server-entry";
import { forwardApiRequest, forwardHealthRequest } from "./server/elysia-adapter";

export default createServerEntry({
  fetch(request) {
    const pathname = new URL(request.url).pathname;

    if (pathname === "/health") return forwardHealthRequest(request);
    if (pathname === "/api" || pathname.startsWith("/api/")) return forwardApiRequest(request);

    return handler.fetch(request);
  },
});
```

- [ ] **Step 4: Run the adapter tests.**

Run: `rtk bun test apps/web/src/server/elysia-adapter.test.ts`

Expected: PASS. Task 3 generates the route tree after it creates the required root route.

- [ ] **Step 5: Commit the embedded adapter.**

```bash
rtk git add apps/web/src/server/elysia-adapter.ts apps/web/src/server/elysia-adapter.test.ts apps/web/src/server.ts apps/web/src/server.test.ts
rtk git commit -m "feat(web): embed Elysia through TanStack Start server entry"
```

### Task 3: Migrate the application shell and page routes

**Files:**
- Create: `apps/web/src/router.tsx`
- Create: `apps/web/src/routes/__root.tsx`
- Create: `apps/web/src/routes/index.tsx`
- Create: `apps/web/src/routes/dashboard.tsx`
- Create: `apps/web/src/components/query-provider.tsx`
- Create: `apps/web/src/styles.css`
- Delete: `apps/web/app/layout.tsx`
- Delete: `apps/web/app/(public)/page.tsx`
- Delete: `apps/web/app/(dashboard)/dashboard/page.tsx`
- Delete: `apps/web/app/globals.css`
- Delete: `apps/web/components/query-provider.tsx`

- [ ] **Step 1: Write a failing route-tree assertion.**

Create `apps/web/src/router.test.ts`:

```ts
import { expect, test } from "bun:test";

test("registers public and dashboard routes", async () => {
  const { getRouter } = await import("./router");
  const router = getRouter();

  expect(router.routeTree.children).toHaveProperty("/");
  expect(router.routeTree.children).toHaveProperty("/dashboard");
});
```

- [ ] **Step 2: Run the test to verify it fails because the TanStack router does not exist.**

Run: `rtk bun test apps/web/src/router.test.ts`

Expected: FAIL with `Cannot find module './router'`.

- [ ] **Step 3: Implement the router, root document, providers, and page routes.**

Create `apps/web/src/router.tsx`:

```tsx
import { createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";

export function getRouter() {
  return createRouter({
    routeTree,
    defaultPreload: "intent",
    scrollRestoration: true,
  });
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>;
  }
}
```

Create `apps/web/src/components/query-provider.tsx`:

```tsx
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState, type ReactNode } from "react";

export function QueryProvider({ children }: { children: ReactNode }) {
  const [queryClient] = useState(() => new QueryClient());

  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
}
```

Create `apps/web/src/routes/__root.tsx`:

```tsx
import { ColorSchemeScript, createTheme, MantineProvider } from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import { HeadContent, Scripts, createRootRoute } from "@tanstack/react-router";
import type { ReactNode } from "react";
import { QueryProvider } from "@/components/query-provider";
import appCss from "@/styles.css?url";

const theme = createTheme({
  primaryColor: "oriskin",
  defaultRadius: "md",
  fontFamily: '"Space Grotesk", system-ui, sans-serif',
  headings: { fontFamily: '"Space Grotesk", system-ui, sans-serif' },
  colors: {
    oriskin: ["#fff1f0", "#ffe1df", "#ffc4bf", "#ff9d96", "#f87168", "#e84f45", "#d93a30", "#b82d25", "#982821", "#7e251f"],
  },
});

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: "Application Workspace" },
      { name: "description", content: "A reusable TanStack Start application workspace" },
    ],
    links: [
      { rel: "preconnect", href: "https://fonts.googleapis.com" },
      { rel: "preconnect", href: "https://fonts.gstatic.com", crossOrigin: "anonymous" },
      { rel: "stylesheet", href: "https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@300..700&display=swap" },
      { rel: "stylesheet", href: appCss },
    ],
  }),
  shellComponent: RootDocument,
});

function RootDocument({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <head>
        <ColorSchemeScript defaultColorScheme="light" />
        <HeadContent />
      </head>
      <body>
        <MantineProvider defaultColorScheme="light" theme={theme}>
          <QueryProvider>
            <Notifications position="top-right" />
            {children}
          </QueryProvider>
        </MantineProvider>
        <Scripts />
      </body>
    </html>
  );
}
```

Create `apps/web/src/routes/index.tsx` and `apps/web/src/routes/dashboard.tsx` with the current page content, replacing `next/link` with TanStack's `Link`:

```tsx
// src/routes/index.tsx
import { Link, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({ component: PublicPage });

function PublicPage() {
  return (
    <main>
      <h1>Application workspace</h1>
      <p>Public entry point for your next product.</p>
      <Link to="/dashboard">Open dashboard</Link>
    </main>
  );
}

// src/routes/dashboard.tsx
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/dashboard")({ component: DashboardPage });

function DashboardPage() {
  return (
    <main>
      <h1>Dashboard</h1>
      <p>Authenticated access can be added by the product built from this template.</p>
    </main>
  );
}
```

Move the existing contents of `apps/web/app/globals.css` into `apps/web/src/styles.css`, delete the old App Router files and old provider, then regenerate `src/routeTree.gen.ts`.

- [ ] **Step 4: Run the route test, typecheck, and production build.**

Run: `rtk bun test apps/web/src/router.test.ts && rtk bun run --cwd apps/web check-types && rtk bun run --cwd apps/web build`

Expected: PASS; Vite generates a route tree and Nitro emits `apps/web/.output/server/index.mjs`.

- [ ] **Step 5: Commit the TanStack application shell.**

```bash
rtk git add apps/web/src
rtk git add -u apps/web/app apps/web/components/query-provider.tsx
rtk git commit -m "refactor(web): migrate pages to TanStack Router"
```

### Task 4: Preserve Eden boundaries and remove Next public configuration

**Files:**
- Create: `apps/web/src/lib/api/client.ts`
- Create: `apps/web/src/lib/api/client.test.ts`
- Create: `apps/web/src/lib/api/server.ts`
- Delete: `apps/web/lib/api/client.ts`
- Delete: `apps/web/lib/api/client.test.ts`
- Delete: `apps/web/lib/api/server.ts`
- Delete: `apps/web/lib/env.ts`
- Modify: `config/web.yaml`
- Modify: `.env.example`
- Modify: `docker-compose.yml`
- Modify: `packages/config/src/run.ts`
- Modify: `packages/config/src/loader.test.ts`

- [ ] **Step 1: Write failing tests for a same-origin browser client and Vite public filtering.**

Create `apps/web/src/lib/api/client.test.ts`:

```ts
import { expect, test } from "bun:test";

test("builds an Eden client at the same-origin API base path", async () => {
  const { api, createApiClient } = await import("./client");

  expect(api).toBeDefined();
  expect(createApiClient("/api")).toBeDefined();
});
```

In `packages/config/src/loader.test.ts`, replace the `loadConfig` import with:

```ts
import { publicConfigEnvironment } from "./run";
import { loadConfig } from "./loader";
```

Then add this test:

```ts
test("exposes only explicit Vite public values", () => {
  withConfigDirectory(
    { web: "VITE_APP_TITLE: Demo\nDATABASE_URL: ${DATABASE_URL}\n" },
    (configDirectory) => {
      expect(
        publicConfigEnvironment(
          ["web"],
          { DATABASE_URL: "postgresql://secret" },
          configDirectory,
        ),
      ).toEqual({ VITE_APP_TITLE: "Demo" });
    },
  );
});
```

- [ ] **Step 2: Run the tests to verify the client and Vite filtering are not present.**

Run: `rtk bun test apps/web/src/lib/api/client.test.ts packages/config/src/loader.test.ts`

Expected: FAIL because the new client file is absent and `publicConfigEnvironment` still filters `NEXT_PUBLIC_` keys.

- [ ] **Step 3: Implement the Eden facades and remove `NEXT_PUBLIC_APP_URL`.**

Create `apps/web/src/lib/api/client.ts`:

```ts
import { treaty } from "@elysia/eden";
import type { App } from "@repo/api";

type ApiClient = ReturnType<typeof treaty<App>>["api"];

export function createApiClient(baseUrl: string): ApiClient {
  return treaty<App>(baseUrl).api;
}

export const api: ApiClient = createApiClient("/api");
```

Create `apps/web/src/lib/api/server.ts`:

```ts
import { treaty } from "@elysia/eden";
import { app } from "@repo/api";
import type { App } from "@repo/api";

type ServerApi = ReturnType<typeof treaty<App>>["api"];

export const api: ServerApi = treaty<App>(app).api;
```

Change `publicConfigEnvironment` in `packages/config/src/run.ts` to filter `key.startsWith("VITE_")`. Remove `NEXT_PUBLIC_APP_URL` from `config/web.yaml`, `.env.example`, Compose build args and environments, and the former Next test. Do not introduce a `VITE_*` variable: the baseline browser client has a constant same-origin API base path.

- [ ] **Step 4: Run the focused client/config tests and search for forbidden public-config references.**

Run: `rtk bun test apps/web/src/lib/api/client.test.ts packages/config/src/loader.test.ts`

Expected: tests PASS.

Then run: `rtk rg -n "NEXT_PUBLIC_APP_URL|@t3-oss/env-nextjs|server-only" apps/web config .env.example docker-compose.yml packages/config`

Expected: `rg` exits 1 with no matches.

- [ ] **Step 5: Commit the configuration and Eden migration.**

```bash
rtk git add apps/web/src/lib/api config/web.yaml .env.example docker-compose.yml packages/config/src/run.ts packages/config/src/loader.test.ts
rtk git add -u apps/web/lib
rtk git commit -m "refactor(web): use same-origin Eden with Vite config"
```

### Task 5: Ship the Bun/Nitro web image and Compose service

**Files:**
- Modify: `Dockerfile.web`
- Modify: `docker-compose.yml`
- Modify: `scripts/dockerfiles.test.ts`

- [ ] **Step 1: Add failing assertions for the generated web output.**

Append to `scripts/dockerfiles.test.ts`:

```ts
test("Dockerfile.web runs the Bun-targeted TanStack Start output", () => {
  const dockerfile = readFileSync(join(root, "Dockerfile.web"), "utf8");

  expect(dockerfile).toContain("bun run vite build");
  expect(dockerfile).toContain("/app/apps/web/.output");
  expect(dockerfile).toContain(".output/server/index.mjs");
  expect(dockerfile).not.toContain(".next");
  expect(dockerfile).not.toContain("next start");
});
```

- [ ] **Step 2: Run the test to verify the current image is Next-specific.**

Run: `rtk bun test scripts/dockerfiles.test.ts`

Expected: FAIL because `Dockerfile.web` contains `.next` and `next build`.

- [ ] **Step 3: Replace the web Docker runtime commands and assets.**

In the builder stage, remove `NEXT_TELEMETRY_DISABLED`, `NEXT_PUBLIC_APP_URL`, `NEXT_PRIVATE_SKIP_PATCHING`, and their arguments. Keep `DATABASE_URL` for Prisma generation, copy `config/`, and build with:

```dockerfile
RUN cd apps/web && bun run ../../packages/config/src/run.ts --modules base,web,api,storage -- bun run vite build
```

In the runner stage, copy the complete generated output and configuration, then run the Bun server:

```dockerfile
COPY --from=builder --chown=nodejs:nodejs /app/apps/web/.output ./.output
COPY --from=builder --chown=nodejs:nodejs /app/packages/config ./packages/config
COPY --from=builder --chown=nodejs:nodejs /app/config ./config
USER nodejs
EXPOSE 4100
CMD ["bun", "run", "packages/config/src/run.ts", "--modules", "base,web,api,storage", "--", "bun", ".output/server/index.mjs"]
```

Retain the existing Turbo prune and dependency stages. The `ssr.noExternal` list from Task 1 keeps the workspace API, configuration, database, and storage packages in the server bundle, so the runner does not need a full source-tree copy.

Remove the obsolete `NEXT_PUBLIC_APP_URL` Compose build argument and environment values. Keep `PORT=4100`, `DATABASE_URL`, the port mapping, and the health check.

- [ ] **Step 4: Run the Dockerfile test and build the web image.**

Run: `rtk bun test scripts/dockerfiles.test.ts && rtk docker build -f Dockerfile.web -t vkit-rapid-web:tanstack-start .`

Expected: PASS; Docker finishes with a Bun image whose command is `.output/server/index.mjs`.

- [ ] **Step 5: Start the Compose web service and verify both embedded endpoints.**

Run: `rtk docker compose up --build -d web && rtk docker compose exec web bun -e "await fetch('http://127.0.0.1:4100/health').then((response) => { if (!response.ok) throw new Error(String(response.status)); })" && rtk docker compose exec web bun -e "await fetch('http://127.0.0.1:4100/api/status').then((response) => { if (!response.ok) throw new Error(String(response.status)); })"`

Expected: all commands PASS. Stop the service afterwards with `rtk docker compose down`.

- [ ] **Step 6: Commit the deployment migration.**

```bash
rtk git add Dockerfile.web docker-compose.yml scripts/dockerfiles.test.ts apps/web/vite.config.ts
rtk git commit -m "build(web): run TanStack Start with Bun and Nitro"
```

### Task 6: Update reusable documentation and complete verification

**Files:**
- Modify: `README.md`
- Modify: `.agent/web.md`
- Modify: `.agent/architecture.md`
- Modify: `.agent/config.md`
- Modify: `docs/superpowers/specs/2026-07-18-tanstack-start-migration-design.md` only if implementation changed an explicit decision

- [ ] **Step 1: Write failing documentation assertions.**

Create `apps/web/docs.test.ts`:

```ts
import { expect, test } from "bun:test";

test("documents TanStack Start as the web runtime", async () => {
  const readme = await Bun.file(new URL("../../README.md", import.meta.url)).text();
  const webRules = await Bun.file(new URL("../../.agent/web.md", import.meta.url)).text();

  expect(readme).toContain("TanStack Start");
  expect(readme).not.toContain("Next.js for the web experience");
  expect(webRules).toContain("TanStack Start");
  expect(webRules).toContain("src/server.ts");
  expect(webRules).not.toContain("App Router");
});
```

- [ ] **Step 2: Run the test to confirm the old Next.js documentation is still present.**

Run: `rtk bun test apps/web/docs.test.ts`

Expected: FAIL because the README and web rules still name Next.js and App Router.

- [ ] **Step 3: Update the public architecture and agent rules.**

Replace all architecture-level legacy-framework references with TanStack Start, TanStack Router, Vite, and Bun/Nitro terminology. Document the invariant that `src/server.ts` is the only TanStack server entry that imports Elysia and intercepts `/api/*` and `/health`; all browser requests use Eden at same-origin `/api`; and browser-visible configuration must be explicitly named `VITE_*`. Keep the existing API/usecase/database/queue/runtime ownership language unchanged.

Update README development, build, Docker, and architecture examples so they state that `task dev` runs TanStack Start with embedded Elysia on port 4100, while `task dev:standalone-api` runs Elysia at port 4101.

- [ ] **Step 4: Run focused and repository-wide verification.**

Run: `rtk task test:web && rtk task check-types:web && rtk task build:web && rtk task quality && rtk task build`

Expected: every command exits 0. Then run `rtk task dev:web` in one terminal and, from another, `rtk task web:health` plus `rtk curl --fail http://localhost:4100/api/status`; both endpoint checks exit 0.

- [ ] **Step 5: Review the final diff and commit documentation.**

Run: `rtk git diff --check && rtk git status --short`

Expected: no whitespace errors; only migration files and any pre-existing user changes appear. Stage only the files listed in this plan and commit:

```bash
rtk git add README.md .agent/web.md .agent/architecture.md .agent/config.md apps/web/docs.test.ts
rtk git commit -m "docs: document TanStack Start web architecture"
```

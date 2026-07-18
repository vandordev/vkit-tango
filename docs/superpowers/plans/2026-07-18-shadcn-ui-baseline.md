# shadcn/ui Baseline Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace Mantine with a source-owned shadcn/ui and Tailwind CSS v4 baseline in the TanStack Start web app while documenting shadcn/ui as the repository default.

**Architecture:** Tailwind CSS v4 is loaded by the existing Vite runtime, and shadcn primitives live in `apps/web/src/components/ui` as ordinary project source. The root route remains responsible for the document and global fallbacks, but removes Mantine providers in favor of semantic HTML and the local `Button`. TanStack Start, Router, Nitro, Eden, and React Query wiring remain unchanged.

**Tech Stack:** Bun workspaces, TanStack Start, TanStack Router, React 19, Vite 8, Nitro, Tailwind CSS v4, shadcn/ui conventions, Radix UI, Lucide React, Bun test, ESLint, TypeScript.

---

## File structure

| File | Responsibility |
| --- | --- |
| `apps/web/package.json` | Replace Mantine dependencies with the minimal Tailwind/shadcn runtime dependencies. |
| `bun.lock` | Record the workspace dependency resolution. |
| `apps/web/vite.config.ts` | Register Tailwind’s Vite plugin alongside the current TanStack and Nitro plugins. |
| `apps/web/components.json` | Define the Vite, non-RSC shadcn project configuration and aliases. |
| `apps/web/src/lib/utils.ts` | Export shadcn’s `cn` class-name merge helper. |
| `apps/web/src/components/ui/button.tsx` | Provide the source-owned Button primitive, including `asChild` support. |
| `apps/web/src/styles.css` | Load Tailwind and define the global shadcn semantic token layer plus retained table utilities. |
| `apps/web/src/routes/__root.tsx` | Render the document and default fallback states without Mantine. |
| `apps/web/src/routes/index.tsx` | Use the local button for the public navigation action. |
| `apps/web/src/routes/dashboard.tsx` | Provide the matching shadcn-based dashboard starter state. |
| `apps/web/vite.config.test.ts` | Assert Tailwind/shadcn configuration and Mantine removal from the package manifest. |
| `apps/web/src/routes/__root.test.ts` | Assert fallback copy, recovery behavior, and shadcn button usage. |
| `apps/web/docs.test.ts` | Assert repository guidance chooses shadcn/ui but permits alternative primary UI systems. |
| `AGENTS.md`, `.agent/ui.md`, `README.md` | Define shadcn/ui as the baseline and document alternatives correctly. |
| Existing `docs/superpowers/specs/*.md` and `docs/superpowers/plans/*.md` listed in Task 4 | Mark prior Mantine-baseline guidance as superseded. |

### Task 1: Establish failing regression tests for the shadcn baseline

**Files:**
- Modify: `apps/web/vite.config.test.ts`
- Modify: `apps/web/src/routes/__root.test.ts`
- Modify: `apps/web/docs.test.ts`

- [ ] **Step 1: Extend the Vite and manifest test with the desired baseline assertions**

  In `apps/web/vite.config.test.ts`, add this test after the existing Vite configuration test:

  ```ts
  test("configures Tailwind and shadcn dependencies without Mantine", async () => {
    const viteConfig = await Bun.file(new URL("./vite.config.ts", import.meta.url)).text();
    const componentsConfig = await Bun.file(new URL("./components.json", import.meta.url)).text();
    const { dependencies } = await Bun.file(new URL("./package.json", import.meta.url)).json();

    expect(viteConfig).toContain('from "@tailwindcss/vite"');
    expect(viteConfig).toContain("tailwindcss()");
    expect(componentsConfig).toContain('"rsc": false');
    expect(componentsConfig).toContain('"css": "src/styles.css"');
    expect(componentsConfig).toContain('"ui": "@/components/ui"');
    expect(dependencies["@mantine/core"]).toBeUndefined();
    expect(dependencies["@mantine/hooks"]).toBeUndefined();
    expect(dependencies["@mantine/notifications"]).toBeUndefined();
    expect(dependencies.tailwindcss).toBeDefined();
    expect(dependencies["class-variance-authority"]).toBeDefined();
    expect(dependencies.clsx).toBeDefined();
    expect(dependencies["radix-ui"]).toBeDefined();
    expect(dependencies["tailwind-merge"]).toBeDefined();
  });
  ```

- [ ] **Step 2: Extend the root fallback source test**

  Add these assertions to the existing test in `apps/web/src/routes/__root.test.ts`:

  ```ts
  expect(rootRoute).toContain('from "@/components/ui/button"');
  expect(rootRoute).not.toContain("@mantine/");
  expect(rootRoute).not.toContain("MantineProvider");
  expect(rootRoute).not.toContain("ColorSchemeScript");
  expect(rootRoute).not.toContain("Notifications");
  ```

- [ ] **Step 3: Replace the documentation test with the new baseline contract**

  Replace `apps/web/docs.test.ts` with:

  ```ts
  import { expect, test } from "bun:test";

  test("documents shadcn ui as the default and alternatives as project choices", async () => {
    const [readme, agentRules, uiRules] = await Promise.all([
      Bun.file(new URL("../../README.md", import.meta.url)).text(),
      Bun.file(new URL("../../AGENTS.md", import.meta.url)).text(),
      Bun.file(new URL("../../.agent/ui.md", import.meta.url)).text(),
    ]);

    for (const document of [readme, agentRules, uiRules]) {
      expect(document).toContain("shadcn/ui");
      expect(document).toContain("one primary UI system");
      expect(document).toContain("Mantine");
      expect(document).toContain("MUI");
    }

    expect(readme).toContain("TanStack Start");
    expect(readme).not.toContain("Mantine by default");
    expect(agentRules).not.toContain("Mantine is the current default");
  });
  ```

- [ ] **Step 4: Run the focused tests and confirm the expected failure**

  Run:

  ```bash
  rtk bun test apps/web/vite.config.test.ts apps/web/src/routes/__root.test.ts apps/web/docs.test.ts
  ```

  Expected: FAIL because `components.json`, Tailwind dependencies/plugin, shadcn button import, and revised guidance do not exist yet.

- [ ] **Step 5: Commit the failing-test checkpoint**

  ```bash
  rtk git add apps/web/vite.config.test.ts apps/web/src/routes/__root.test.ts apps/web/docs.test.ts
  rtk git commit -m "test(web): specify shadcn ui baseline"
  ```

### Task 2: Install the source-owned Tailwind and shadcn foundation

**Files:**
- Modify: `apps/web/package.json`
- Modify: `bun.lock`
- Modify: `apps/web/vite.config.ts`
- Create: `apps/web/components.json`
- Create: `apps/web/src/lib/utils.ts`
- Create: `apps/web/src/components/ui/button.tsx`
- Modify: `apps/web/src/styles.css`

- [ ] **Step 1: Replace Mantine packages with the minimal shadcn dependency set**

  In `apps/web/package.json`, remove `@mantine/core`, `@mantine/hooks`, and
  `@mantine/notifications`. Add these `dependencies`, retaining all existing
  non-Mantine entries and their versions:

  ```json
  {
    "@tailwindcss/vite": "^4.1.18",
    "class-variance-authority": "^0.7.1",
    "clsx": "^2.1.1",
    "radix-ui": "^1.6.2",
    "tailwind-merge": "^3.0.2",
    "tailwindcss": "^4.1.18"
  }
  ```

  Then update the lockfile without lifecycle scripts:

  ```bash
  rtk bun install --ignore-scripts
  ```

- [ ] **Step 2: Register Tailwind in the existing Vite configuration**

  Add the import:

  ```ts
  import tailwindcss from "@tailwindcss/vite";
  ```

  Add `tailwindcss()` immediately before `tanstackRouter(...)` in the `plugins`
  array of `apps/web/vite.config.ts`. Retain both router test-ignore settings,
  Nitro’s Bun preset, React, the `@` alias, and `ssr.noExternal` unchanged.

- [ ] **Step 3: Add the shadcn project configuration and utility helper**

  Create `apps/web/components.json` with:

  ```json
  {
    "$schema": "https://ui.shadcn.com/schema.json",
    "style": "new-york",
    "rsc": false,
    "tsx": true,
    "tailwind": {
      "config": "",
      "css": "src/styles.css",
      "baseColor": "zinc",
      "cssVariables": true,
      "prefix": ""
    },
    "aliases": {
      "components": "@/components",
      "utils": "@/lib/utils",
      "ui": "@/components/ui",
      "lib": "@/lib"
    },
    "iconLibrary": "lucide"
  }
  ```

  Create `apps/web/src/lib/utils.ts` with:

  ```ts
  import type { ClassValue } from "clsx";
  import { clsx } from "clsx";
  import { twMerge } from "tailwind-merge";

  export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
  }
  ```

- [ ] **Step 4: Add the local shadcn Button primitive**

  Create `apps/web/src/components/ui/button.tsx` with:

  ```tsx
  import type * as React from "react";
  import { cva, type VariantProps } from "class-variance-authority";
  import { Slot } from "radix-ui";

  import { cn } from "@/lib/utils";

  const buttonVariants = cva(
    "inline-flex h-9 shrink-0 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium whitespace-nowrap transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
    {
      variants: {
        variant: {
          default: "bg-primary text-primary-foreground hover:bg-primary/90",
          destructive: "bg-destructive text-destructive-foreground hover:bg-destructive/90",
          outline: "border border-input bg-background hover:bg-accent hover:text-accent-foreground",
          secondary: "bg-secondary text-secondary-foreground hover:bg-secondary/80",
          ghost: "hover:bg-accent hover:text-accent-foreground",
          link: "text-primary underline-offset-4 hover:underline",
        },
        size: {
          default: "h-9 px-4 py-2",
          sm: "h-8 rounded-md px-3 text-xs",
          lg: "h-10 rounded-md px-8",
          icon: "size-9 p-0",
        },
      },
      defaultVariants: { variant: "default", size: "default" },
    },
  );

  function Button({
    className,
    variant,
    size,
    asChild = false,
    ...props
  }: React.ComponentProps<"button"> &
    VariantProps<typeof buttonVariants> & { asChild?: boolean }) {
    const Comp = asChild ? Slot.Root : "button";

    return <Comp className={cn(buttonVariants({ variant, size, className }))} {...props} />;
  }

  export { Button, buttonVariants };
  ```

- [ ] **Step 5: Replace the global Mantine stylesheet with the Tailwind token layer**

  Replace the Mantine imports at the top of `apps/web/src/styles.css` with:

  ```css
  @import "tailwindcss";

  @theme inline {
    --color-background: var(--background);
    --color-foreground: var(--foreground);
    --color-primary: var(--primary);
    --color-primary-foreground: var(--primary-foreground);
    --color-secondary: var(--secondary);
    --color-secondary-foreground: var(--secondary-foreground);
    --color-muted: var(--muted);
    --color-muted-foreground: var(--muted-foreground);
    --color-accent: var(--accent);
    --color-accent-foreground: var(--accent-foreground);
    --color-destructive: var(--destructive);
    --color-destructive-foreground: var(--destructive-foreground);
    --color-border: var(--border);
    --color-input: var(--input);
    --color-ring: var(--ring);
    --radius-sm: calc(var(--radius) - 4px);
    --radius-md: calc(var(--radius) - 2px);
    --radius-lg: var(--radius);
  }

  :root {
    --background: oklch(0.985 0.003 285.8);
    --foreground: oklch(0.18 0.01 285.8);
    --primary: oklch(0.33 0.14 24);
    --primary-foreground: oklch(0.985 0.003 285.8);
    --secondary: oklch(0.95 0.008 25);
    --secondary-foreground: oklch(0.26 0.02 25);
    --muted: oklch(0.95 0.006 285.8);
    --muted-foreground: oklch(0.48 0.015 285.8);
    --accent: oklch(0.93 0.012 25);
    --accent-foreground: oklch(0.26 0.02 25);
    --destructive: oklch(0.52 0.19 25);
    --destructive-foreground: oklch(0.985 0.003 285.8);
    --border: oklch(0.89 0.008 285.8);
    --input: oklch(0.89 0.008 285.8);
    --ring: oklch(0.48 0.16 24);
    --radius: 0.5rem;
  }
  ```

  Keep the existing global reset and table utility selectors, but map their custom
  colors to the semantic variables (`--background`, `--foreground`, `--muted-
  foreground`, and `--border`) so existing table consumers continue to render.
  Remove every `@mantine` import and all `--app-*` declarations.

- [ ] **Step 6: Run the focused tests and typecheck the web workspace**

  Run:

  ```bash
  rtk bun test apps/web/vite.config.test.ts
  rtk task check-types:web
  ```

  Expected: the Vite/manifest test passes and TypeScript accepts the local button
  foundation; root and documentation tests still fail until the next tasks.

- [ ] **Step 7: Commit the foundation**

  ```bash
  rtk git add apps/web/package.json bun.lock apps/web/vite.config.ts apps/web/components.json apps/web/src/lib/utils.ts apps/web/src/components/ui/button.tsx apps/web/src/styles.css
  rtk git commit -m "feat(web): add shadcn ui foundation"
  ```

### Task 3: Rebuild the starter pages and fallback states using shadcn source

**Files:**
- Modify: `apps/web/src/routes/__root.tsx`
- Modify: `apps/web/src/routes/index.tsx`
- Modify: `apps/web/src/routes/dashboard.tsx`
- Modify: `apps/web/src/routes/__root.test.ts`

- [ ] **Step 1: Remove Mantine from the root route while retaining its routing behavior**

  In `apps/web/src/routes/__root.tsx`:

  - Remove all `@mantine/*` imports, `ColorSchemeScript`, `createTheme`, the
    `theme` constant, and `Notifications`.
  - Import `Link`, `HeadContent`, `Scripts`, `createRootRoute`, and `useNavigate`
    from `@tanstack/react-router`; retain `ErrorComponentProps` and `ReactNode`
    as type-only imports.
  - Import `Button` from `@/components/ui/button`, `QueryProvider`, and `appCss`.
  - Keep the existing route metadata, fallback registration, `import.meta.env.DEV`
    error detail guard, `reset`, and `navigate({ to: "/" })` actions.
  - Use this fallback layout:

  ```tsx
  function FallbackLayout({ children }: { children: ReactNode }) {
    return <main className="mx-auto flex min-h-screen w-full max-w-xl flex-col justify-center gap-6 px-6 py-24">{children}</main>;
  }
  ```

  - Render headings as `h1` with `text-3xl font-semibold tracking-tight`, muted
    text as `p` with `text-muted-foreground`, and action rows as `div` with
    `flex flex-wrap gap-3`.
  - Render the home actions as `Button asChild` wrapping `<Link to="/">Back to
    home</Link>`. Keep `Try again` as `<Button onClick={reset}>Try again</Button>`.
  - Render the document body as:

  ```tsx
  <body className="min-h-screen bg-background font-sans text-foreground antialiased">
    <QueryProvider>{children}</QueryProvider>
    <Scripts />
  </body>
  ```

- [ ] **Step 2: Use the local Button in both starter routes**

  Replace `apps/web/src/routes/index.tsx` with:

  ```tsx
  import { Link, createFileRoute } from "@tanstack/react-router";

  import { Button } from "@/components/ui/button";

  export const Route = createFileRoute("/")({ component: PublicPage });

  function PublicPage() {
    return (
      <main className="mx-auto flex min-h-screen w-full max-w-3xl flex-col justify-center gap-6 px-6 py-24">
        <div className="space-y-2">
          <h1 className="text-3xl font-semibold tracking-tight">Application workspace</h1>
          <p className="text-muted-foreground">Public entry point for your next product.</p>
        </div>
        <div>
          <Button asChild>
            <Link to="/dashboard">Open dashboard</Link>
          </Button>
        </div>
      </main>
    );
  }
  ```

  Replace `apps/web/src/routes/dashboard.tsx` with:

  ```tsx
  import { Link, createFileRoute } from "@tanstack/react-router";

  import { Button } from "@/components/ui/button";

  export const Route = createFileRoute("/dashboard")({ component: DashboardPage });

  function DashboardPage() {
    return (
      <main className="mx-auto flex min-h-screen w-full max-w-3xl flex-col justify-center gap-6 px-6 py-24">
        <div className="space-y-2">
          <h1 className="text-3xl font-semibold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground">Authenticated access can be added by the product built from this template.</p>
        </div>
        <div>
          <Button asChild variant="outline">
            <Link to="/">Back to home</Link>
          </Button>
        </div>
      </main>
    );
  }
  ```

- [ ] **Step 3: Make the fallback source test assert actual shadcn navigation use**

  Add these assertions to `apps/web/src/routes/__root.test.ts`:

  ```ts
  expect(rootRoute).toContain("<Button asChild>");
  expect(rootRoute).toContain('<Link to="/">Back to home</Link>');
  expect(rootRoute).toContain("bg-background");
  ```

- [ ] **Step 4: Run the focused fallback and routing tests**

  Run:

  ```bash
  rtk bun test apps/web/src/routes/__root.test.ts apps/web/src/router.test.ts
  rtk task check-types:web
  ```

  Expected: PASS.

- [ ] **Step 5: Commit the page migration**

  ```bash
  rtk git add apps/web/src/routes/__root.tsx apps/web/src/routes/index.tsx apps/web/src/routes/dashboard.tsx apps/web/src/routes/__root.test.ts
  rtk git commit -m "refactor(web): replace mantine pages with shadcn"
  ```

### Task 4: Align all guidance with the new baseline and verify the repository

**Files:**
- Modify: `AGENTS.md`
- Modify: `.agent/ui.md`
- Modify: `README.md`
- Modify: `docs/superpowers/specs/2026-07-18-tanstack-start-migration-design.md`
- Modify: `docs/superpowers/specs/2026-07-18-tanstack-start-default-states-design.md`
- Modify: `docs/superpowers/specs/2026-07-12-reusable-elysia-eden-boilerplate-design.md`
- Modify: `docs/superpowers/plans/2026-07-18-tanstack-start-migration.md`
- Modify: `docs/superpowers/plans/2026-07-18-tanstack-start-default-states.md`
- Modify: `docs/superpowers/plans/2026-07-12-reusable-elysia-eden-foundation.md`
- Modify: `apps/web/docs.test.ts`

- [ ] **Step 1: Update primary repository guidance**

  Make these exact policy changes:

  - In `AGENTS.md`, replace the repository-shape claim with “UI choice is
    project-scoped: shadcn/ui is the current baseline; Mantine, MUI, and other
    component libraries are valid alternatives, not second baselines.” Replace
    the UI rule with “Use one primary UI system per project. Do not mix shadcn,
    Mantine, MUI, or other primitive libraries as defaults.”
  - Replace `.agent/ui.md` with four bullets: shadcn/ui is the current source-owned
    Tailwind baseline; a project may deliberately choose Mantine, MUI, or another
    component library; it must choose one primary UI system; a switch removes
    prior providers, dependencies, theme code, and primitives rather than keeping
    competing systems active.
  - In `README.md`, describe the baseline as “shadcn/ui by default, with Mantine,
    MUI, or another library as a deliberate alternative”; change the current-web
    paragraph and the `apps/web` table cell from Mantine to shadcn/ui.

- [ ] **Step 2: Mark former Mantine documentation as superseded history**

  In the two 2026-07-18 specs and plans, replace claims that Mantine remains the
  baseline or that the fallback uses Mantine with a direct note that those details
  describe the implementation at the time and are superseded by
  `2026-07-18-shadcn-ui-baseline-design.md`. Update tech-stack/table references
  from Mantine to `shadcn/ui and Tailwind CSS v4`; do not retain example Mantine
  imports as current instructions.

  In the 2026-07-12 historical design and plan, preserve the statement that the
  original boilerplate used Mantine only as past tense history, and append a short
  “Superseded UI baseline” note that names shadcn/ui as the current default and
  points to the new design specification.

- [ ] **Step 3: Run a repository-wide guidance scan**

  Run:

  ```bash
  rtk rg -n -i 'Mantine (is|remains|by default|current default)|current web baseline uses Mantine|Mantine UI' AGENTS.md README.md .agent docs --glob '*.md'
  ```

  Expected: no matches. Mentions of Mantine must only describe a permitted
  alternative or clearly superseded historical state.

- [ ] **Step 4: Run all verification commands**

  Run:

  ```bash
  rtk bun test apps/web/docs.test.ts
  rtk task quality
  rtk task build
  rtk git diff --check
  rtk git status --short
  ```

  Expected: all tests, linting, typechecks, and builds pass; whitespace check is
  clean; only the intended migration files remain staged or uncommitted.

- [ ] **Step 5: Commit documentation and final verification state**

  ```bash
  rtk git add AGENTS.md .agent/ui.md README.md docs/superpowers/specs/2026-07-18-tanstack-start-migration-design.md docs/superpowers/specs/2026-07-18-tanstack-start-default-states-design.md docs/superpowers/specs/2026-07-12-reusable-elysia-eden-boilerplate-design.md docs/superpowers/plans/2026-07-18-tanstack-start-migration.md docs/superpowers/plans/2026-07-18-tanstack-start-default-states.md docs/superpowers/plans/2026-07-12-reusable-elysia-eden-foundation.md apps/web/docs.test.ts
  rtk git commit -m "docs: make shadcn ui the baseline"
  rtk git status --short
  ```

  Expected: `git status --short` prints no output.

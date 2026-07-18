# TanStack Start Default States Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add English, production-safe global not-found and error fallback states to the TanStack Start web application.

**Architecture:** Keep fallback ownership in `src/routes/__root.tsx`, where TanStack Router resolves root-level `notFoundComponent` and `errorComponent`. The current components use source-owned shadcn/ui primitives and TanStack Router links; no API or new route is needed. The Mantine implementation details below are historical and superseded by [the shadcn/ui baseline design](../specs/2026-07-18-shadcn-ui-baseline-design.md).

**Tech Stack:** TanStack Router, TanStack Start, React 19, shadcn/ui, Tailwind CSS v4, Bun test.

---

### Task 1: Register and render global fallback states

**Files:**
- Modify: `apps/web/src/routes/__root.tsx`
- Modify: `apps/web/src/router.test.ts`

- [ ] **Step 1: Write the failing root-route test.**

Add this assertion to `apps/web/src/router.test.ts`:

```ts
test("registers global not-found and error fallback states", async () => {
  const { Route } = await import("./routes/__root");

  expect(Route.options.notFoundComponent).toBeDefined();
  expect(Route.options.errorComponent).toBeDefined();
});
```

- [ ] **Step 2: Run the test and confirm it fails.**

Run: `rtk bun test apps/web/src/router.test.ts`

Expected: FAIL because the root route does not configure fallback components.

- [ ] **Step 3: Add the fallback components.**

The historical Mantine import instructions in this task are superseded. In the current `apps/web/src/routes/__root.tsx`, import the local shadcn `Button`, `Link`, `type ErrorComponentProps`, and `createRootRoute`; register these options on the existing root route:

```tsx
notFoundComponent: NotFoundPage,
errorComponent: RouteErrorPage,
```

Implement a shared `FallbackLayout` that returns a centered `<main>` containing a `Stack`. Implement `NotFoundPage` with `Title` text `Page not found`, body text `The page you are looking for does not exist or has moved.`, and a `Button component={Link} to="/"` labeled `Back to home`.

Implement `RouteErrorPage({ error, reset }: ErrorComponentProps)` with `Title` text `Something went wrong`, generic body text `Please try again. If the problem continues, return to the home page.`, a primary `Button onClick={reset}` labeled `Try again`, and an outline home button. Render `error.message` only when `import.meta.env.DEV` is true and `error` is an `Error`.

- [ ] **Step 4: Run the focused tests.**

Run: `rtk bun test apps/web/src/router.test.ts`

Expected: PASS with both route-registration tests green.

- [ ] **Step 5: Commit the fallback behavior.**

```bash
rtk git add apps/web/src/routes/__root.tsx apps/web/src/router.test.ts
rtk git commit -m "feat(web): add TanStack Start fallback states"
```

### Task 2: Verify the fallback copy and recovery controls

**Files:**
- Create: `apps/web/src/routes/__root.test.ts`
- Modify: `apps/web/src/routes/__root.tsx`

- [ ] **Step 1: Write the failing copy and recovery-control test.**

Create `apps/web/src/routes/__root.test.ts`:

```ts
import { expect, test } from "bun:test";

test("provides English not-found and error recovery actions", async () => {
  const rootRoute = await Bun.file(new URL("./__root.tsx", import.meta.url)).text();

  expect(rootRoute).toContain("Page not found");
  expect(rootRoute).toContain("Back to home");
  expect(rootRoute).toContain("Something went wrong");
  expect(rootRoute).toContain("Try again");
  expect(rootRoute).toContain("onClick={reset}");
  expect(rootRoute).toContain("import.meta.env.DEV");
});
```

- [ ] **Step 2: Run the test and confirm it fails.**

Run: `rtk bun test apps/web/src/routes/__root.test.ts`

Expected: FAIL because the fallback copy and reset control do not exist.

- [ ] **Step 3: Keep recovery controls accessible.**

Ensure the action row wraps on narrow viewports. Use shadcn semantic tokens and avoid custom gradients, decorative illustrations, or new CSS files.

- [ ] **Step 4: Run focused tests and typecheck.**

Run: `rtk bun test apps/web/src/router.test.ts apps/web/src/routes/__root.test.ts && rtk task check-types:web`

Expected: all tests and typecheck pass.

- [ ] **Step 5: Commit tests and polish.**

```bash
rtk git add apps/web/src/routes/__root.tsx apps/web/src/routes/__root.test.ts apps/web/src/router.test.ts
rtk git commit -m "test(web): cover fallback recovery states"
```

### Task 3: Verify the complete web runtime

**Files:**
- Verify: `apps/web/src/routes/__root.tsx`
- Verify: `apps/web/src/routes/__root.test.ts`
- Verify: `apps/web/src/router.test.ts`

- [ ] **Step 1: Run repository verification.**

Run: `rtk task quality && rtk task build`

Expected: every test, lint task, typecheck, and build task exits 0.

- [ ] **Step 2: Confirm the generated route tree remains stable.**

Run: `rtk git diff --check && rtk git status --short`

Expected: no generated route-tree drift and only intentional committed changes remain.

# shadcn/ui Baseline Design

**Date:** 2026-07-18

## Goal

Make `shadcn/ui` the active UI baseline for the TanStack Start web application. The
repository must no longer depend on Mantine or describe it as the default. This is
a starter convention, not a restriction: a downstream project may choose Mantine,
MUI, or another UI system, provided it adopts one primary system rather than
mixing competing primitive libraries by default.

## Scope

### Web runtime and component foundation

- Replace Mantine with Tailwind CSS v4 and the shadcn/ui Vite setup used by the
  sibling TanStack Start reference project.
- Add `components.json` for a non-RSC TanStack Start application. It will point
  to `apps/web/src/styles.css`, use CSS variables, Lucide icons, and the existing
  TypeScript aliases.
- Add the minimal shadcn support packages and a source-owned `cn` helper in
  `apps/web/src/lib/utils.ts`.
- Add source-owned shadcn primitives under `apps/web/src/components/ui/`. The
  initial baseline needs `Button`; future primitives are copied into this folder
  as the application needs them, rather than being wrapped behind a second local
  design system.
- Preserve the current router, API, React Query, and Nitro wiring. This is a UI
  migration only.

### Styling and page states

- Replace Mantine stylesheet imports, provider, theme object, color-scheme script,
  and notifications provider with Tailwind’s CSS entrypoint and shadcn semantic
  tokens in `src/styles.css`.
- Keep a modest neutral starter appearance with accessible foreground, border,
  focus-ring, destructive, and radius tokens. Do not introduce an unused theme
  switcher or notification library.
- Rebuild the root document, not-found screen, route-error screen, home page, and
  dashboard using semantic HTML, Tailwind utilities, and the local shadcn button.
  The existing English fallback copy and recovery actions remain unchanged:
  `Page not found`, `Back to home`, `Something went wrong`, and `Try again`.
- Use `Button asChild` with TanStack Router links where a navigation control needs
  shadcn styling. Keep route navigation and error reset behavior intact.

### Removal and documentation

- Remove every Mantine package from `apps/web/package.json` and all Mantine imports
  or runtime references from web source and styles.
- Update `AGENTS.md`, `.agent/ui.md`, and `README.md` to say that shadcn/ui is the
  repository baseline, while Mantine, MUI, and other libraries remain valid
  project-level alternatives. The guidance must still require one primary UI
  system per project.
- Update the active migration and fallback-state specs/plans, plus any remaining
  Markdown guidance that presents Mantine as the current baseline. Historical
  documents may retain factual implementation history only when they clearly mark
  it as superseded; no Markdown document may recommend Mantine as this
  repository’s default.

## Dependency and configuration decisions

The implementation will use the same compatible foundation as
`../my-tanstack-start`:

- `tailwindcss` and `@tailwindcss/vite` for Tailwind CSS v4.
- `class-variance-authority`, `clsx`, and `tailwind-merge` for the standard
  shadcn button and `cn` helper.
- `radix-ui` for the shadcn `Slot` used by `Button asChild`.

The Vite configuration adds the Tailwind plugin without disturbing TanStack
Router, TanStack Start, Nitro, or React plugin ordering. The repository’s existing
`@/*` aliases remain the canonical aliases used by `components.json`.

## Tests and verification

- Update source-level tests to assert the shadcn/Tailwind configuration and the
  continued root fallback behavior.
- Add regression checks that the web package and source do not reference Mantine,
  and that the repository guidance names shadcn/ui as the baseline while allowing
  alternative primary UI systems.
- Run the focused tests first, then `task quality` and `task build`.

## Non-goals

- Redesigning the product or adding a demo business domain.
- Adding every shadcn primitive, dark-mode controls, toast infrastructure, or
  third-party component registries before a feature needs them.
- Preventing consuming projects from using Mantine, MUI, or another component
  library.

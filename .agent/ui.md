# UI Rules

- UI libraries are project-level choices, not backend architecture requirements.
- The current web baseline is shadcn/ui: Tailwind-based, source-owned primitives live in `apps/web/src/components/ui`.
- A project may deliberately choose Mantine, MUI, or another component library instead of shadcn/ui.
- Choose one primary UI system per project. Do not keep shadcn/ui, Mantine, MUI, or other primitive libraries as competing defaults.
- When changing systems, remove the previous providers, dependencies, theme code, and primitives instead of leaving competing systems active.
- Keep UI components inside `apps/web`; do not create `packages/ui` until multiple apps genuinely share the same implementation.
- Shared UI primitives belong in `apps/web/src/components/ui`; feature-specific components stay close to their route or feature.
- Preserve the selected system's accessibility, responsive behavior, loading states, error states, and interaction conventions.

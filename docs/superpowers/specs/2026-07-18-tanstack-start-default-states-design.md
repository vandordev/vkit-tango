# TanStack Start Default States Design

## Goal

Provide production-ready default not-found and route-error states for the TanStack Start web application. All user-visible copy is English.

## Scope

The global root route owns both fallback states:

- A not-found state handles unregistered URLs.
- An error boundary handles errors thrown by a route component, loader, or other routed UI.

Both states are lightweight Mantine layouts that preserve the application's existing warm, restrained visual language. They use the existing shared shell and do not add a new design system, route, API request, or dependency.

## Components and behavior

`apps/web/src/routes/__root.tsx` will register `notFoundComponent` and `errorComponent` on `createRootRoute`.

The not-found component presents a clear “Page not found” message and a single “Back to home” action using TanStack Router navigation.

The error component presents “Something went wrong” with two actions:

- “Try again” calls the boundary reset callback.
- “Back to home” uses TanStack Router navigation.

The production message stays generic. In development, the error message is shown below the generic explanation to make local debugging practical without exposing technical details in production.

## Visual direction

The states use a centered, compact content column with ordinary typography, clear spacing, and standard Mantine buttons. They use the existing coral primary color only for the primary action. There are no decorative illustrations, gradients, floating cards, or framework branding.

## Testing and acceptance criteria

Focused tests will verify that the root route exposes both fallback handlers and that their English labels and recovery actions are available. Existing route registration, typecheck, lint, and build checks must remain green.

The feature is accepted when an unknown URL renders the not-found state, a route error renders the error boundary, and both states provide a working route back to `/`.

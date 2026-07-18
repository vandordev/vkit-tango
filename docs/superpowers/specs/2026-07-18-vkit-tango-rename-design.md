# vkit-tango Project Rename Design

**Status:** Approved
**Date:** 2026-07-18

## Goal

Rename the active project identity from `vkit-fast` to `vkit-tango` without
changing runtime behavior, database schema, or historical design records.

## Active identity boundary

The rename applies to all active, user-facing and build-facing identifiers:

- the README title and Go API OpenAPI title;
- the `app.name` YAML value and Docker Compose project name;
- the root Bun workspace package name;
- the Go module path and every Go import, including committed Ent generated
  output; and
- the project-name regression test.

The Go module becomes `github.com/vandordev/vkit-tango`. This is necessary so
the source tree, generated Ent imports, Go tooling, and published project
identity describe the same project.

## Explicit exclusions

The generic local database name `boilerplate` remains unchanged. It is a
development data-store default, not public project identity. Historical
documents under `docs/superpowers/specs/` and `docs/superpowers/plans/` retain
their original wording to preserve an accurate implementation record.

## Verification

The identity regression test is changed first to require `vkit-tango`, then
the active identifiers are changed until that test passes. `go test ./...`,
`task quality`, `task build`, and the Ent/River PostgreSQL integration test run
against a temporary PostgreSQL 16 container complete the verification.

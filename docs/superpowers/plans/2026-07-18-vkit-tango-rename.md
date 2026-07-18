# vkit-tango Project Rename Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rename every active project identity from `vkit-fast` to `vkit-tango` while preserving runtime behavior and historical documentation.

**Architecture:** The root Go module becomes `github.com/vandordev/vkit-tango`; all Go source and committed Ent output import that canonical module. Runtime-facing labels change in configuration, Compose, OpenAPI, README, and root Bun metadata. Historical specs and plans remain immutable records.

**Tech Stack:** Go 1.25, Ent generated Go code, Bun, Docker Compose, Bun test, Task.

---

### Task 1: Establish a failing active-identity contract

**Files:**
- Modify: `scripts/project-name.test.ts`
- Test: `scripts/project-name.test.ts`

- [ ] **Step 1: Change the expected public and Compose identifiers.**

Replace the test assertions with:

```ts
test("uses vkit-tango for public and Compose identities", () => {
  const readme = await Bun.file("README.md").text();
  const compose = await Bun.file("docker-compose.yml").text();
  const appConfig = await Bun.file("config/app.yaml").text();

  expect(readme).toContain("# vkit-tango");
  expect(compose).toMatch(/^name: vkit-tango$/m);
  expect(appConfig).toContain("name: vkit-tango");
});
```

- [ ] **Step 2: Verify the regression test fails.**

Run: `rtk bun test scripts/project-name.test.ts`

Expected: FAIL because the active files still contain `vkit-fast`.

### Task 2: Rename active identity and Go module references

**Files:**
- Modify: `README.md`
- Modify: `config/app.yaml`
- Modify: `docker-compose.yml`
- Modify: `package.json`
- Modify: `go.mod`
- Modify: `apps/**/*.go`, `internal/**/*.go`, `tools/**/*.go`
- Modify: `contracts/openapi/openapi.json`
- Modify: `scripts/project-name.test.ts`
- Exclude: `docs/superpowers/specs/**`, `docs/superpowers/plans/**`

- [ ] **Step 1: Replace active text identities.**

Apply these exact active values:

```text
README heading: # vkit-tango
config/app.yaml app.name: vkit-tango
docker-compose.yml top-level name: vkit-tango
package.json name: vkit-tango
go.mod module: github.com/vandordev/vkit-tango
OpenAPI info.title: vkit-tango API
```

- [ ] **Step 2: Rewrite all active Go module imports.**

Replace every `github.com/vandordev/vkit-fast/` prefix outside historical docs
with `github.com/vandordev/vkit-tango/`, including `internal/platform/db/**`
generated Ent files. Do not alter package declarations or file paths.

- [ ] **Step 3: Verify the focused identity test passes.**

Run: `rtk bun test scripts/project-name.test.ts`

Expected: PASS with `vkit-tango` in README, Compose, and app configuration.

- [ ] **Step 4: Verify no active legacy identity remains.**

Run:

```bash
rtk rg -n 'vkit-fast|github.com/vandordev/vkit-fast|turbostack' --glob '!docs/superpowers/specs/**' --glob '!docs/superpowers/plans/**' --glob '!node_modules/**' .
```

Expected: no matches.

### Task 3: Verify generated contracts and build outputs

**Files:**
- Modify if generated: `contracts/openapi/openapi.json`
- Modify if generated: `apps/web/src/lib/api/generated/**`

- [ ] **Step 1: Regenerate and check the OpenAPI-driven web client.**

Run: `rtk task api:client:check`

Expected: exit 0 and no diff in `contracts/openapi` or generated web client.

- [ ] **Step 2: Run full repository verification.**

Run:

```bash
rtk task quality
rtk task build
rtk git diff --check
```

Expected: all commands exit 0; generated TypeScript lint warnings may remain warnings but must not be errors.

- [ ] **Step 3: Run the Ent/River integration proof against temporary PostgreSQL.**

Start PostgreSQL with:

```bash
rtk docker run --detach --rm --name vkit-tango-integration-postgres -e POSTGRES_DB=vkit_test -e POSTGRES_USER=vkit_test -e POSTGRES_PASSWORD=vkit_test -p 127.0.0.1::5432 postgres:16-alpine
```

Read the mapped port, wait for `pg_isready`, then run:

```bash
rtk env 'TEST_DATABASE_URL=postgresql://vkit_test:vkit_test@127.0.0.1:<port>/vkit_test?sslmode=disable' go test ./internal/platform/river -run TestEntAndRiverShareTransaction -count=1
rtk docker stop vkit-tango-integration-postgres
```

Expected: the integration test passes and the temporary container is removed.

- [ ] **Step 4: Commit the active rename.**

```bash
git add README.md config/app.yaml docker-compose.yml package.json go.mod apps internal tools contracts/openapi scripts/project-name.test.ts
git commit -m "chore: rename project to vkit tango"
```

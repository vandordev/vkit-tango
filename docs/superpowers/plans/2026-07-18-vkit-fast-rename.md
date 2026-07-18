# vkit-fast Rename Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rename the current boilerplate identity from `vkit-rapid` to `vkit-fast` without changing the adjacent `../vkit-rapid` project.

**Architecture:** Update the public repository name and Docker Compose project name, while preserving historical migration documents as records of the old implementation. Add a focused regression test that protects user-facing naming and Compose resource naming.

**Tech Stack:** Bun test, Docker Compose, Markdown, YAML.

---

### Task 1: Protect the new public identity

**Files:**
- Create: `scripts/project-name.test.ts`
- Modify: `README.md`
- Modify: `docker-compose.yml`

- [ ] **Step 1: Write the failing test.**

```ts
test("uses vkit-fast for public and Compose identities", () => {
  expect(readme).toContain("# vkit-fast");
  expect(readme).not.toContain("vkit-rapid");
  expect(compose).toMatch(/^name: vkit-fast$/m);
});
```

- [ ] **Step 2: Run the test and confirm it fails.**

Run: `rtk bun test scripts/project-name.test.ts`

Expected: FAIL because the README still names `vkit-rapid` and Compose has no explicit project name.

- [ ] **Step 3: Apply the minimal rename.**

Replace the README title, descriptive prose, and quick-start `cd` command with `vkit-fast`. Add `name: vkit-fast` as the top-level Compose project name so service, network, and volume names do not inherit the workspace directory name.

- [ ] **Step 4: Run focused verification.**

Run: `rtk bun test scripts/project-name.test.ts && rtk docker compose config --quiet`

Expected: PASS with valid Compose configuration.

- [ ] **Step 5: Commit.**

```bash
rtk git add README.md docker-compose.yml scripts/project-name.test.ts docs/superpowers/plans/2026-07-18-vkit-fast-rename.md
rtk git commit -m "chore: rename boilerplate to vkit-fast"
```

### Task 2: Verify no active identity remains

**Files:**
- Verify: `README.md`
- Verify: `docker-compose.yml`
- Verify: `scripts/project-name.test.ts`

- [ ] **Step 1: Run repository verification.**

Run: `rtk task quality && rtk task build && rtk rg -n -i "vkit-rapid" README.md docker-compose.yml scripts --glob '!node_modules/**'`

Expected: quality and build pass; the final search exits 1 with no active `vkit-rapid` identity.

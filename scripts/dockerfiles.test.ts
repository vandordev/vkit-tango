import { readFileSync } from "node:fs";
import { join } from "node:path";
import { expect, test } from "bun:test";

const root = join(import.meta.dir, "..");

for (const [file, workspace] of [["Dockerfile.api", "@repo/api"], ["Dockerfile.worker", "@repo/worker"], ["Dockerfile.scheduler", "@repo/scheduler"]] as const) {
	test(`${file} prunes and ships only production dependencies`, () => {
		const dockerfile = readFileSync(join(root, file), "utf8");
		expect(dockerfile).toContain("FROM oven/bun:1.3.14-alpine AS base");
		expect(dockerfile).toContain(`turbo@2.10.5 prune ${workspace} --docker`);
    expect(dockerfile).toContain("bun install --production --frozen-lockfile --ignore-scripts");
    expect(dockerfile).not.toContain("COPY --from=builder /app/node_modules ./node_modules");
  });
}

test("Dockerfile.web builds from a pruned graph", () => {
	expect(readFileSync(join(root, "Dockerfile.web"), "utf8")).toContain("turbo@2.10.5 prune web --docker");
});

test("Dockerfile.web runs the Bun-targeted TanStack Start output", () => {
	const dockerfile = readFileSync(join(root, "Dockerfile.web"), "utf8");

	expect(dockerfile).toContain("FROM oven/bun:1.3.14-alpine AS base");
	expect(dockerfile).toContain("bun run vite build");
	expect(dockerfile).toContain("/app/apps/web/.output");
	expect(dockerfile).toContain(".output/server/index.mjs");
	expect(dockerfile).toContain("COPY --from=deps /app/apps/web/node_modules ./apps/web/node_modules");
	expect(dockerfile).toContain("COPY --from=deps /app/packages/config/node_modules ./packages/config/node_modules");
	expect(dockerfile).toContain("RUN apps/web/node_modules/.bin/prisma generate --schema=packages/database/prisma/schema.prisma");
	expect(dockerfile).toContain("RUN mkdir -p /runtime/config-node-modules && cp -LR packages/config/node_modules/. /runtime/config-node-modules/");
	expect(dockerfile).toContain("RUN rm -rf ./packages/config/node_modules");
	expect(dockerfile).toContain("COPY --from=builder --chown=app:app /runtime/config-node-modules ./packages/config/node_modules");
	expect(dockerfile).not.toContain(".next");
	expect(dockerfile).not.toContain("next start");
	expect(dockerfile).not.toContain("nextjs");
});

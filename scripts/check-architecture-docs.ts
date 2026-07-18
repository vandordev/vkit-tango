const files = ["README.md", "AGENTS.md", ".agent/architecture.md", ".agent/api.md", ".agent/database.md", ".agent/config.md", ".agent/worker.md", ".agent/web.md"];
const required = ["Huma", "Ent", "Goose", "River", "/api/v1", "Hey API", "Socket.IO"];
const requiredByFile: Record<string, string[]> = {
  "README.md": ["Fx", "internal/contract", "task sync", "task dev"],
  "AGENTS.md": ["Fx", "internal/contract", "task add:*", "task sync"],
  ".agent/architecture.md": ["Fx", "internal/contract", "internal/generated/fx"],
  ".agent/api.md": ["Fx", "humachi", "contract.Command", "task sync"],
  ".agent/worker.md": ["Fx", "internal/contract.Command", "apps/scheduler"],
  ".agent/config.md": ["Fx", "scheduler"],
  ".agent/web.md": ["TanStack Start", "Hey API", "/api/v1"],
};
const forbidden = ["Elysia", "Prisma", "pg-boss", "Eden", "standalone scheduler", "schedules run inside the Go worker", "register periodic reconciliation jobs on every replica"];

for (const file of files) {
  const content = await Bun.file(file).text();
  for (const term of required) if (!content.includes(term)) throw new Error(`${file} is missing ${term}`);
  for (const term of requiredByFile[file] ?? []) if (!content.includes(term)) throw new Error(`${file} is missing ${term}`);
  for (const term of forbidden) if (content.includes(term)) throw new Error(`${file} still contains ${term}`);
}

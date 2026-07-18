const files = ["README.md", "AGENTS.md", ".agent/architecture.md", ".agent/api.md", ".agent/database.md", ".agent/config.md", ".agent/worker.md"];
const required = ["Huma", "Ent", "Goose", "River", "/api/v1", "Hey API", "Socket.IO"];
const forbidden = ["Elysia", "Prisma", "pg-boss", "Eden", "standalone scheduler"];

for (const file of files) {
  const content = await Bun.file(file).text();
  for (const term of required) if (!content.includes(term)) throw new Error(`${file} is missing ${term}`);
  for (const term of forbidden) if (content.includes(term)) throw new Error(`${file} still contains ${term}`);
}

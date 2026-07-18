const output = await new Response(Bun.spawn(["task", "--list"], { stdout: "pipe" }).stdout).text();

for (const task of [
  "dev:api",
  "dev:worker",
	"dev:scheduler",
  "dev:migrate",
  "api:openapi",
  "api:client:generate",
  "api:client:check",
  "db:generate",
  "db:migrate",
  "db:migrate:status",
]) {
  if (!output.includes(task)) throw new Error(`Missing required task: ${task}`);
}

for (const task of ["start:scheduler", "db:studio"]) {
  if (output.includes(task)) throw new Error(`Obsolete task is still present: ${task}`);
}

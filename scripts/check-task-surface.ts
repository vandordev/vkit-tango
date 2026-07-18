const output = await new Response(Bun.spawn(["task", "--list"], { stdout: "pipe" }).stdout).text();

for (const task of [
  "dev",
  "dev:api",
  "dev:worker",
	"dev:scheduler",
  "dev:migrate",
  "ci",
  "doctor",
  "test:go",
  "test:web",
  "test:realtime",
  "test:config",
  "api:openapi",
  "api:client:generate",
  "api:client:check",
  "db:generate",
  "db:migrate",
  "db:migrate:status",
]) {
  if (!output.includes(task)) throw new Error(`Missing required task: ${task}`);
}

const taskfile = await Bun.file("Taskfile.yml").text();
if (!taskfile.includes("CLI_ARGS_LIST")) throw new Error("dev task must forward selected services from CLI arguments");
if (!taskfile.includes("dev:api dev:worker dev:scheduler dev:web")) throw new Error("dev task must define the default service set");
if (!taskfile.includes('cmds: ["rtk task quality", "rtk task build"]')) throw new Error("ci task must run quality before build");
if (!taskfile.includes("for tool in go bun task vx docker") || !taskfile.includes("docker compose version") || !taskfile.includes("missing .env")) throw new Error("doctor task must check required local setup");
const syncBlock = taskfile.slice(taskfile.indexOf("  sync:\n"), taskfile.indexOf("  test:\n"));
if (!syncBlock.includes("rtk task api:client:generate")) throw new Error("sync task must refresh OpenAPI and Hey API output");

for (const task of ["start:scheduler", "db:studio"]) {
  if (output.includes(task)) throw new Error(`Obsolete task is still present: ${task}`);
}

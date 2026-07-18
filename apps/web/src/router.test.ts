import { expect, test } from "bun:test";

test("registers public and dashboard routes", async () => {
  const { getRouter } = await import("./router");
  const router = getRouter();

  expect(router.routesByPath).toHaveProperty("/");
  expect(router.routesByPath).toHaveProperty("/dashboard");
});

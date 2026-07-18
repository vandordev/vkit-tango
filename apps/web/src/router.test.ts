import { expect, test } from "bun:test";

test("registers public and dashboard routes", async () => {
  const { getRouter } = await import("./router");
  const router = getRouter();

  expect(router.routesByPath).toHaveProperty("/");
  expect(router.routesByPath).toHaveProperty("/dashboard");
});

test("registers global not-found and error fallback states", async () => {
  const { Route } = await import("./routes/__root");

  expect(Route.options.notFoundComponent).toBeDefined();
  expect(Route.options.errorComponent).toBeDefined();
});

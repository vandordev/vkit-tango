import { expect, test } from "bun:test";

import { createQueryClient } from "./query-provider";

test("invalidates cached queries after a successful mutation", async () => {
  const queryClient = createQueryClient();
  const queryKey = ["system-metadata"];
  queryClient.setQueryData(queryKey, { key: "maintenance" });

  const mutation = queryClient.getMutationCache().build(queryClient, {
    mutationFn: async () => ({ ok: true }),
  });
  await mutation.execute(undefined);

  expect(queryClient.getQueryState(queryKey)?.isInvalidated).toBe(true);
});

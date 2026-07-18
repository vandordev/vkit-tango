import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "../../contracts/openapi/openapi.json",
  output: {
    path: "src/lib/api/generated",
    postProcess: ["prettier"],
  },
  plugins: ["@hey-api/client-fetch", "@hey-api/sdk", "@tanstack/react-query"],
});

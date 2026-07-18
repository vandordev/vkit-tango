import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import { defineConfig } from "vite";
import { nitro } from "nitro/vite";

export default defineConfig({
  plugins: [
    tailwindcss(),
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
      routeFileIgnorePattern: "\\.test\\.",
    }),
    tanstackStart({ router: { routeFileIgnorePattern: "\\.test\\." } }),
    nitro({ preset: "bun" }),
    react(),
  ],
  resolve: {
    alias: {
      "@": new URL("./src", import.meta.url).pathname,
    },
  },
  ssr: {
    noExternal: ["@repo/config"],
  },
  server: {
    proxy: {
      "/api": { target: "http://localhost:4101", changeOrigin: true },
      "/health": { target: "http://localhost:4101", changeOrigin: true },
    },
  },
});

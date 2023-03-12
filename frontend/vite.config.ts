import {defineConfig} from "vite";
import {svelte} from "@sveltejs/vite-plugin-svelte";
import sveltePreprocess from "svelte-preprocess";
import {resolve as pathResolve} from "path";
import basicSsl from "@vitejs/plugin-basic-ssl";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    https: true,
  },
  plugins: [
    svelte({
      preprocess: sveltePreprocess({
        preserve: ["ld+json"],
        scss: {
          includePaths: ["src/"],
          quietDeps: true,
        },
      }),
    }),
    basicSsl(),
  ],
  optimizeDeps: {exclude: ["svelte-navigator"]},
  resolve: {
    alias: {
      "~": pathResolve(__dirname, "src"),
    },
  },
});

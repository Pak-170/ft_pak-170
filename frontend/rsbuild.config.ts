import { defineConfig } from '@rsbuild/core';
import { pluginSvelte } from '@rsbuild/plugin-svelte';
import { pluginSass } from '@rsbuild/plugin-sass';

export default defineConfig({
  plugins: [
    pluginSvelte({
      svelteLoaderOptions: {},
    }),
    pluginSass(),
  ],
  server: {
    proxy: {
      '/api': 'http://localhost:3001',
    },
  },
});

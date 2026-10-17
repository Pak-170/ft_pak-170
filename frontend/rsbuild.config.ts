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
      '/api': {
        target: 'http://backend:8080',
        ws: true,
      },
    },
  },
});

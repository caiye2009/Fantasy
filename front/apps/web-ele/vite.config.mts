import { defineConfig } from '@vben/vite-config';
import { fileURLToPath, URL } from 'node:url';

import ElementPlus from 'unplugin-element-plus/vite';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      resolve: {
        alias: {
          '@': fileURLToPath(new URL('./src', import.meta.url)),
        },
      },
      plugins: [
        ElementPlus({
          format: 'esm',
        }),
      ],
      server: {
        proxy: {
          '/api': {
            target: 'http://localhost:8081',
            changeOrigin: true,
            rewrite: (path) => path,
            ws: true,
          },
        },
      },
    },
  };
});

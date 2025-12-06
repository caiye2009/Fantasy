import { defineConfig } from '@vben/vite-config';

import ElementPlus from 'unplugin-element-plus/vite';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      plugins: [
        ElementPlus({
          format: 'esm',
        }),
      ],
      server: {
        proxy: {
          '/api': {
            target: 'http://localhost:8081', // 改为您的后端地址
            changeOrigin: true,
            rewrite: (path) => path, // 保持 /api 前缀
            ws: true,
          },
        },
      },
    },
  };
});
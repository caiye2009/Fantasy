import { defineConfig } from '@vben/vite-config';
import ElementPlus from 'unplugin-element-plus/vite';
import AutoImport from 'unplugin-auto-import/vite';
import Components from 'unplugin-vue-components/vite';
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      plugins: [
        // 自动导入 Element Plus 相关函数，如：ElMessage, ElMessageBox
        AutoImport({
          resolvers: [ElementPlusResolver()],
          dts: 'auto-imports.d.ts',
        }),
        // 自动导入 Element Plus 组件
        Components({
          resolvers: [ElementPlusResolver()],
          dts: 'components.d.ts',
        }),
        // Element Plus 样式按需导入
        ElementPlus({
          format: 'esm',
        }),
      ],
      server: {
        port: 5777,
        proxy: {
          '/api/v1': {
            target: 'http://localhost:8081',
            changeOrigin: true,
            ws: true,
            rewrite: (path) => path.replace(/^\/api\/v1/, '/api/v1'),
          },
        },
      },
    },
  };
});
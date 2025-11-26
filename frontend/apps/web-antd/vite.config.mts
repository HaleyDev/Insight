import { defineConfig } from '@vben/vite-config';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      server: {
        proxy: {
          '/api': {
            changeOrigin: true,
            // 不要重写路径，保持 /api 前缀
            // rewrite: (path) => path.replace(/^\/api/, ''),
            // 后端服务器地址，不要包含 /api 路径
            target: 'http://10.99.65.75:8099',
            ws: true,
          },
        },
      },
    },
  };
});

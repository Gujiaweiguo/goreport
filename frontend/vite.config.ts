import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, __dirname, '')
  const proxyTarget = env.VITE_PROXY_TARGET || 'http://localhost:8085'

  return {
    root: path.resolve(__dirname),
    plugins: [vue()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src')
      }
    },
    build: {
      target: 'esnext',
      minify: 'esbuild',
      sourcemap: mode === 'development',
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (!id.includes('node_modules')) {
              return
            }

            if (id.includes('node_modules/zrender')) {
              return 'zrender'
            }

            if (id.includes('node_modules/echarts')) {
              return 'echarts'
            }

            if (id.includes('node_modules/element-plus')) {
              return 'element-plus'
            }

            if (
              id.includes('node_modules/vue') ||
              id.includes('node_modules/vue-router') ||
              id.includes('node_modules/pinia')
            ) {
              return 'vendor'
            }
          }
        }
      },
      chunkSizeWarningLimit: 1000,
      assetsInlineLimit: 4096
    },
    server: {
      port: 3000,
      proxy: {
        '/health': {
          target: proxyTarget,
          changeOrigin: true
        },
        '/api': {
          target: proxyTarget,
          changeOrigin: true
        },
        '/datasource': {
          target: proxyTarget,
          changeOrigin: true
        },
        '/jmreport': {
          target: proxyTarget,
          changeOrigin: true
        }
      }
    }
  }
})

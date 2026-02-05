import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  root: path.resolve(__dirname),
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/health': {
        target: 'http://localhost:8085',
        changeOrigin: true
      },
      '/api': {
        target: 'http://localhost:8085',
        changeOrigin: true
      },
      '/datasource': {
        target: 'http://localhost:8085',
        changeOrigin: true
      },
      '/jmreport': {
        target: 'http://localhost:8085',
        changeOrigin: true
      }
    }
  }
})

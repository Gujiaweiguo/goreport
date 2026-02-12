import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  cacheDir: '.vitest-cache',
  test: {
    globals: true,
    environment: 'jsdom',
    environmentOptions: {
      url: 'http://localhost:3000'
    },
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/tests/',
        '**/*.d.ts',
        '**/*.test.ts',
        '**/*.test.tsx',
        '**/main.ts'
      ]
    },
    setupFiles: ['./src/tests/setup.ts']
  }
})

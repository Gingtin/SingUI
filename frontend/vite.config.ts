import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    outDir: '../backend/cmd/server/dist',
    emptyOutDir: true,
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:2096',
        changeOrigin: true,
        ws: true,
      },
      '/sub': {
        target: 'http://127.0.0.1:2096',
        changeOrigin: true,
      },
    },
  },
})

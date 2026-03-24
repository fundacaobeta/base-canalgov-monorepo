import { fileURLToPath, URL } from 'node:url'
import autoprefixer from 'autoprefixer'
import tailwind from 'tailwindcss'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const apiTarget = process.env.VITE_API_TARGET || 'http://127.0.0.1:9000'
const wsTarget  = process.env.VITE_WS_TARGET  || 'ws://127.0.0.1:9000'

export default defineConfig({
  css: {
    postcss: {
      plugins: [tailwind(), autoprefixer()],
    },
  },
  server: {
    port: 8001,
    proxy: {
      '/api':      { target: apiTarget, changeOrigin: true },
      '/uploads':  { target: apiTarget, changeOrigin: true },
      '/widget/ws': { target: wsTarget, ws: true, changeOrigin: true },
    },
  },
  build: {
    outDir: 'dist',
    chunkSizeWarningLimit: 600,
  },
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@widget': fileURLToPath(new URL('./src', import.meta.url)),
      '@shared-ui': fileURLToPath(new URL('./src/shared-ui', import.meta.url)),
      '@main': fileURLToPath(new URL('./src', import.meta.url)),
    }
  },
})

import { fileURLToPath, URL } from 'node:url'
import autoprefixer from 'autoprefixer'
import tailwind from 'tailwindcss'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { vueGrab } from 'vue-grab/plugins/vite'

// Allow overriding backend target via env var for Docker dev setup.
// In Docker: VITE_API_TARGET=http://backend:9000
// Locally:   defaults to http://127.0.0.1:9000
const apiTarget = process.env.VITE_API_TARGET || 'http://127.0.0.1:9000'
const wsTarget  = process.env.VITE_WS_TARGET  || 'ws://127.0.0.1:9000'

export default defineConfig({
  css: {
    postcss: {
      plugins: [tailwind(), autoprefixer()],
    },
  },
  server: {
    port: 8000,
    proxy: {
      '/api':     { target: apiTarget, changeOrigin: true },
      '/logout':  { target: apiTarget, changeOrigin: true },
      '/uploads': { target: apiTarget, changeOrigin: true },
      '/ws':      { target: wsTarget,  ws: true, changeOrigin: true },
    },
  },
  build: {
    chunkSizeWarningLimit: 600,
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          'radix': ['radix-vue', 'reka-ui'],
          'icons': ['lucide-vue-next', '@radix-icons/vue'],
          'utils': ['@vueuse/core', 'clsx', 'tailwind-merge', 'class-variance-authority'],
          'charts': ['@unovis/ts', '@unovis/vue'],
          'editor': ['@tiptap/vue-3', '@tiptap/starter-kit', '@tiptap/extension-image', '@tiptap/extension-link', '@tiptap/extension-placeholder', '@tiptap/extension-table', '@tiptap/extension-table-cell', '@tiptap/extension-table-header', '@tiptap/extension-table-row'],
          'forms': ['vee-validate', '@vee-validate/zod', 'zod'],
          'table': ['@tanstack/vue-table'],
          'misc': ['axios', 'date-fns', 'mitt', 'qs', 'vue-i18n']
        }
      }
    }
  },
  plugins: [
    vue(),
    vueGrab(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
})

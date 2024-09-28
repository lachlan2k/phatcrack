import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },

  server: {
    proxy: {
      '/api': {
        target: process.env.DEV_API_URL,
        changeOrigin: true,
        ws: true
      },
      '/agent-assets': {
        target: process.env.DEV_AGENT_SERVER_URL,
        changeOrigin: true
      }
    }
  }
})

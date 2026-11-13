import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig(({ mode }) => {
  // Base path aplikasi. '/' = di root domain (default, perilaku lama).
  // Isi VITE_BASE_PATH=/coffee/ bila Coffee disajikan di bawah path /coffee/ (production).
  const env = loadEnv(mode, process.cwd(), 'VITE_')
  const base = env.VITE_BASE_PATH || '/'

  return {
    base,
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      port: 3002,
      proxy: {
        '/api': {
          target: 'http://localhost:8083',
          changeOrigin: true,
        },
      },
    },
    build: {
      outDir: 'dist',
      sourcemap: false,
      chunkSizeWarningLimit: 1000,
      rollupOptions: {
        output: {
          manualChunks(id: string) {
            if (id.includes('react-dom') || id.includes('react-router-dom') || id.includes('node_modules/react/')) return 'vendor'
            if (id.includes('@tanstack/react-query')) return 'query'
            if (id.includes('@radix-ui')) return 'ui'
            if (id.includes('react-hook-form') || id.includes('zod') || id.includes('@hookform')) return 'form'
            if (id.includes('recharts')) return 'charts'
          },
        },
      },
    },
  }
})

import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  build: {
    // Production build optimizations
    target: 'es2020',
    sourcemap: true, // Enable source maps for error tracking
    rollupOptions: {
      output: {
        // Code splitting for optimal caching
        manualChunks: {
          'vendor-react': ['react', 'react-dom'],
          'vendor-apollo': ['@apollo/client', 'graphql'],
        },
      },
    },
    // Chunk size warning threshold (kB)
    chunkSizeWarningLimit: 500,
  },
  server: {
    // Development proxy to avoid CORS issues
    proxy: {
      '/graphql': {
        target: 'http://localhost:4000',
        changeOrigin: true,
      },
    },
  },
})

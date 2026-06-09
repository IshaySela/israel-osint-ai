import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig(({ mode, command }) => {
  const env = loadEnv(mode, process.cwd())

  if (command === 'build' && !env.VITE_BACKEND_URL) {
    throw new Error('VITE_BACKEND_URL must be defined for production builds')
  }

  return {
    plugins: [
      react(),
      tailwindcss(),
    ],
    server: {
      watch: {
        usePolling: true,
      },
      host: true,
      strictPort: true,
      port: 5173,
    }
  }
})


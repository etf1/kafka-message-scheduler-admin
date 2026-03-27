/// <reference types="vitest" />
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import path from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    // Mirror CRA's baseUrl: "src" — all absolute imports resolve from src/
    alias: {
      _common: path.resolve(__dirname, 'src/_common'),
      _core: path.resolve(__dirname, 'src/_core'),
      business: path.resolve(__dirname, 'src/business'),
      mocks: path.resolve(__dirname, 'src/mocks'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      // Forward /api and /configuration.json to the Go backend
      '/api': { target: 'http://localhost:9000', changeOrigin: true },
      '/configuration.json': {
        target: 'http://localhost:9000',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'build',
    sourcemap: false,
    rollupOptions: {
      output: {
        manualChunks: {
          react: ['react', 'react-dom', 'react-router-dom'],
          i18n: ['i18next', 'react-i18next', 'i18next-browser-languagedetector', 'i18next-http-backend'],
          vendor: ['rxjs', 'date-fns', 'clsx', 'react-transition-group'],
        },
      },
    },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/setupTests.ts'],
  },
});

import { fileURLToPath, URL } from 'url';

import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { VitePWA } from 'vite-plugin-pwa';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    VitePWA({
      registerType: 'autoUpdate',
      workbox: {
        globPatterns: [
          '**/*.{css,html,png,svg,woff,woff2}',
          'manifest.webmanifest',
          '**/index*.js',
          '**/shaka*.js',
          '**/[A-Z]*.js',
        ],
      },
      manifest: {
        name: 'Hyperpipe',
        short_name: 'Hyperpipe',
        start_url: '/',
        display: 'standalone',
        background_color: '#000',
        theme_color: '#000',
        description: 'Privacy respecting YouTube Music Frontend.',
        icons: [
          {
            src: '/android-chrome-192x192.png',
            sizes: '192x192',
            type: 'image/png',
          },
          {
            src: '/android-chrome-512x512.png',
            sizes: '512x512',
            type: 'image/png',
          },
          {
            src: '/android-maskable.png',
            sizes: '1024x1024',
            type: 'image/png',
            purpose: 'maskable',
          },
        ],
      },
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    sourceMap: true,
  },
});

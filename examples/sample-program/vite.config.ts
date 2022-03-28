/**
 * @type {import('vite').UserConfig}
 */
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  build: {
    sourcemap: true,
    manifest: true,
    rollupOptions: {
      input: {
        main: 'src/main.js',
      },
    },
  },
});

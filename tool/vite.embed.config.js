import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  define: {
    'process.env.NODE_ENV': JSON.stringify('production'),
  },
  build: {
    lib: {
      entry: 'src/embed.jsx',
      name: 'BlueprintROI',
      fileName: 'blueprint-roi',
      formats: ['iife'],
    },
    rollupOptions: {
      output: {
        inlineDynamicImports: true,
        entryFileNames: 'blueprint-roi.js',
        assetFileNames: (assetInfo) => {
          if (assetInfo.name === 'style.css') return 'blueprint-roi.css';
          return assetInfo.name;
        },
      },
    },
  },
});

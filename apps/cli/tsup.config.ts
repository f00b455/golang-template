import { defineConfig } from 'tsup';

export default defineConfig({
  entry: ['src/index.ts'],
  format: ['esm'],
  dts: true,
  clean: true,
  shims: true,
  minify: false,
  splitting: false,
  sourcemap: true,
  target: 'node18'
});
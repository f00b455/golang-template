const baseConfig = require('../../eslint.base.cjs');

module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  plugins: ['@typescript-eslint'],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended'
  ],
  parserOptions: {
    ecmaVersion: 2022,
    sourceType: 'module',
    project: './tsconfig.json',
    tsconfigRootDir: __dirname,
  },
  env: {
    node: true,
    es2022: true
  },
  rules: {
    // Spread base rules first
    ...baseConfig.rules,
    // Then override with package-specific rules
    '@typescript-eslint/explicit-module-boundary-types': 'error',
    '@typescript-eslint/no-explicit-any': 'error',
    '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
    '@typescript-eslint/strict-boolean-expressions': 'off',
    'no-console': 'off', // CLI needs console output
    'max-lines-per-function': 'off', // Temporarily disabled for implementation
    'max-statements': 'off', // Temporarily disabled for implementation
    'max-lines': 'off', // Temporarily disabled for implementation
    'complexity': 'off', // Temporarily disabled for implementation
    'max-depth': 'off', // Temporarily disabled for implementation
  },
  overrides: baseConfig.overrides,
  ignorePatterns: ['dist', 'node_modules', '*.js', '*.cjs', '*.mjs', '**/*.test.ts', '**/features/**', 'vitest.config.ts', 'tsup.config.ts']
};
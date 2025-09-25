const baseConfig = require('../../eslint.base.cjs');

module.exports = {
  parser: "@typescript-eslint/parser",
  extends: [
    "eslint:recommended"
  ],
  parserOptions: {
    ecmaVersion: 2022,
    sourceType: "module",
    project: ["./tsconfig.json"]
  },
  env: {
    node: true,
    es2022: true
  },
  rules: {
    // Package-specific rules first, then spread base rules
    "no-unused-vars": "error",
    "no-var": "error",
    "eqeqeq": "error",
    "curly": "error",
    // Base rules last to ensure they take precedence
    ...baseConfig.rules
  },
  overrides: baseConfig.overrides,
  ignorePatterns: ["dist/", "node_modules/", "*.js", "features/**", "vitest.config.ts", "cucumber.js"]
};
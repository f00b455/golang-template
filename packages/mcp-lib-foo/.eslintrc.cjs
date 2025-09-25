const baseConfig = require('../../eslint.base.cjs');

module.exports = {
  root: true,
  parser: "@typescript-eslint/parser",
  plugins: ["@typescript-eslint"],
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended"
  ],
  parserOptions: {
    ecmaVersion: 2022,
    sourceType: "module",
    project: "./tsconfig.json"
  },
  rules: {
    // Package-specific rules first, then spread base rules
    "@typescript-eslint/explicit-module-boundary-types": "error",
    "@typescript-eslint/no-explicit-any": "error",
    "@typescript-eslint/no-unused-vars": ["error", { "argsIgnorePattern": "^_" }],
    "@typescript-eslint/strict-boolean-expressions": "off",
    "no-console": ["warn", { "allow": ["error"] }],
    // Base rules last to ensure they take precedence
    ...baseConfig.rules
  },
  overrides: baseConfig.overrides,
  ignorePatterns: ["dist", "node_modules", "*.js", "*.cjs", "*.mjs"]
};
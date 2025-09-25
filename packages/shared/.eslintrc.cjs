const baseConfig = require('../../eslint.base.cjs');

module.exports = {
  parser: "@typescript-eslint/parser",
  plugins: ["functional"],
  extends: [
    "eslint:recommended",
    "plugin:functional/external-vanilla-recommended",
    "plugin:functional/recommended",
    "plugin:functional/stylistic"
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
    "no-console": "warn",
    "eqeqeq": "error",
    "curly": "error",

    // Pure Function Enforcement
    "functional/no-let": "error",
    "functional/immutable-data": "error",
    "functional/prefer-readonly-type": "error",
    "functional/no-return-void": "error",
    "functional/no-throw-statements": "error",
    "functional/no-expression-statements": ["error", {
      "ignoreVoid": true
    }],

    // Allow some impurity for utility functions
    "functional/no-conditional-statements": "off",
    "functional/no-mixed-types": "off",
    "functional/functional-parameters": "off",
    "functional/readonly-type": "off",

    // Base rules last to ensure they take precedence
    ...baseConfig.rules
  },
  overrides: baseConfig.overrides,
  ignorePatterns: ["dist/", "node_modules/", "*.js", "__tests__/**", "vitest.config.ts", "features/**/*"]
};
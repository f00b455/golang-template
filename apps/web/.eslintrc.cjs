const baseConfig = require('../../eslint.base.cjs');

module.exports = {
  extends: ["next/core-web-vitals"],
  rules: {
    // Package-specific rules first, then spread base rules
    "no-unused-vars": ["error", { "args": "none" }],
    "no-var": "error",
    "no-console": "warn",
    "eqeqeq": "error",
    "curly": "error",
    // Base rules last to ensure they take precedence
    ...baseConfig.rules
  },
  overrides: baseConfig.overrides
};
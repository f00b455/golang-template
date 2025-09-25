module.exports = {
  default: {
    require: ['features/step_definitions/**/*.ts'],
    requireModule: ['tsx'],
    format: ['progress', 'json:test-results/cucumber-report.json'],
    parallel: 1,
    tags: '@pkg-mcp-lib-foo',
  },
};
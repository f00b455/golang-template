module.exports = {
  default: {
    paths: ['features/**/*.feature'],
    require: [
      'features/step_definitions/shared.steps.ts'          // TypeScript ONLY!
    ],
    format: [
      'progress',
      'json:cucumber-report/cucumber-report.json'
    ],
    formatOptions: {},
    parallel: 1,  // ESM safe
    retry: 0,
    strict: true,
  },
};
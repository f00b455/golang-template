module.exports = {
  default: {
    paths: ['features/**/*.feature'],
    require: [
      'features/step_definitions/cli.steps.ts'
    ],
    format: [
      'progress',
      'json:cucumber-report/cucumber-report.json'
    ],
    formatOptions: { snippetInterface: 'async-await' },
  },
};
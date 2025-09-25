export default {
  default: {
    require: [
      'features/step_definitions/**/*.ts'                   // Package-specific steps (TypeScript ONLY!)
    ],
    format: ['progress-bar', 'json:cucumber-report/cucumber_report.json'],
    formatOptions: { snippetInterface: 'async-await' },
  },
};
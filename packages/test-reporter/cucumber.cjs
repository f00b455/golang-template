module.exports = {
  default: {
    paths: ['features/**/*.feature'],
    require: [
      'features/step_definitions/**/*.ts'
    ],
    format: [
      'progress',
      'html:reports/cucumber.html',
      'json:reports/cucumber.json'
    ],
    parallel: 1
  }
};
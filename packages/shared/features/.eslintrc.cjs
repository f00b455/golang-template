module.exports = {
  extends: ["../.eslintrc.cjs"],
  rules: {
    "functional/no-let": "off",
    "functional/no-return-void": "off",
    "functional/no-expression-statements": "off",
    "functional/no-loop-statements": "off",
    "functional/immutable-data": "off",
    "functional/prefer-readonly-type": "off"
  }
};
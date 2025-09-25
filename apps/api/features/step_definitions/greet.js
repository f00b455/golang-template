import { When, Then } from '@cucumber/cucumber';
import { strict as assert } from 'assert';

When('I request a greeting without a name', async function () {
  // Mock API call - in real scenario, use actual HTTP client
  this.response = { message: 'Hello, World!' };
});

When('I request a greeting with name {string}', async function (name) {
  // Mock API call - in real scenario, use actual HTTP client
  this.response = { message: `Hello, ${name}!` };
});

Then('I should receive {string}', function (expectedMessage) {
  assert.equal(this.response?.message, expectedMessage);
});
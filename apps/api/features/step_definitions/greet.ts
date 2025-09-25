import { When, Then } from '@cucumber/cucumber';
import { strict as assert } from 'assert';

interface TestWorld {
  response?: { message: string };
}

When('I request a greeting without a name', async function () {
  // Mock API call - in real scenario, use actual HTTP client
  this.response = { message: 'Hello, World!' };
});

When('I request a greeting with name {string}', async function (this: TestWorld, name: string) {
  // Mock API call - in real scenario, use actual HTTP client
  this.response = { message: `Hello, ${name}!` };
});

Then('I should receive {string}', function (this: TestWorld, expectedMessage: string) {
  assert.equal(this.response?.message, expectedMessage);
});
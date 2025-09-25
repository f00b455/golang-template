import { Given, When, Then } from '@cucumber/cucumber';
import assert from 'assert';
import { greet, formatDate } from '../../dist/src/index.js';

// Module-level context variables (same pattern as lib-foo)
let greetFunction: typeof greet | null = null;
let formatDateFunction: typeof formatDate | null = null;
let inputDate: Date | null = null;
let originalDateString: string | null = null;
let result: any = null;
let results: any[] = [];

// Greet function steps
Given('I have access to the shared greet function', function () {
  greetFunction = greet;
});

When('I call greet with {string}', function (name: string) {
  result = greetFunction!(name);
});

When('I call greet with {string} multiple times', function (name: string) {
  results = [];
  for (let i = 0; i < 5; i++) {
    results.push(greetFunction!(name));
  }
  result = results[0] || '';
});

// Shared steps that were previously in cucumber-shared
Then('the result should be {string}', function (expectedResult: string) {
  assert.equal(result, expectedResult);
});

Then('all results should be identical', function () {
  assert(Array.isArray(results), 'Results should be an array');
  assert(results.length > 1, 'Should have multiple results to compare');
  const firstResult = results[0];
  results.forEach((res, index) => {
    assert.equal(res, firstResult, `Result ${index} should match first result`);
  });
});

// Date formatting function steps
Given('I have access to the shared formatDate function', function () {
  formatDateFunction = formatDate;
});

Given('I have a date {string}', function (dateString: string) {
  inputDate = new Date(dateString);
  originalDateString = inputDate.toISOString();
});

Given('I have a date object', function () {
  inputDate = new Date('2024-06-15T14:30:00Z');
  originalDateString = inputDate.toISOString();
});

When('I call formatDate with that date', function () {
  result = formatDateFunction!(inputDate!);
});

When('I call formatDate with that date multiple times', function () {
  results = [];
  for (let i = 0; i < 5; i++) {
    results.push(formatDateFunction!(inputDate!));
  }
  result = results[0] || '';
});

Then('the original date object should remain unchanged', function () {
  assert.equal(inputDate!.toISOString(), originalDateString, 'Date object should remain unchanged');
});

Then('the function should not modify its input', function () {
  // This is enforced by TypeScript readonly parameter
  // We verify the date hasn't changed
  assert.equal(inputDate!.toISOString(), originalDateString, 'Input date should not be modified');
});
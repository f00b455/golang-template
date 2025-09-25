import { Given, When, Then } from '@cucumber/cucumber';
import { fooTransform, fooFilter, fooProcess, createFooProcessor, fooGreet, FooConfig } from '../../src/index';
import assert from 'assert';

// Data operations context
let inputArray: any[] = [];
let originalArray: any[] = [];
let result: any = null;

// Processing context
let config: FooConfig | null = null;
let processor: ReturnType<typeof createFooProcessor> | null = null;
let processors: ReturnType<typeof createFooProcessor>[] = [];

// Greeting context
let greetFunction: ReturnType<typeof fooGreet> | null = null;

// Data Operations Steps
Given('I have access to the foo data operation functions', function () {
  // Functions are available via imports - no setup needed
  assert(typeof fooTransform === 'function');
  assert(typeof fooFilter === 'function');
});

Given('I have an array {string}', function (arrayString: string) {
  inputArray = JSON.parse(arrayString);
  originalArray = [...inputArray]; // Keep copy of original
});

Given('I have a number array {string}', function (arrayString: string) {
  inputArray = JSON.parse(arrayString);
  originalArray = [...inputArray];
});

Given('I have a string array {string}', function (arrayString: string) {
  inputArray = JSON.parse(arrayString);
  originalArray = [...inputArray];
});

Given('I have an empty array', function () {
  inputArray = [];
  originalArray = [];
});

When('I transform it with uppercase function', function () {
  result = fooTransform(inputArray, (str: string) => str.toUpperCase());
});

When('I transform it with a function that adds {string}', function (suffix: string) {
  result = fooTransform(inputArray, (str: string) => str + suffix);
});

When('I filter it for even numbers', function () {
  result = fooFilter(inputArray, (num: number) => num % 2 === 0);
});

When('I filter it for items containing {string}', function (searchText: string) {
  result = fooFilter(inputArray, (str: string) => str.includes(searchText));
});

When('I transform and filter the array', function () {
  // Transform to uppercase
  const transformed = fooTransform(inputArray, (str: string) => str.toUpperCase());
  // Filter for items longer than 4 characters
  result = fooFilter(transformed, (str: string) => str.length > 4);
});

When('I transform it with any function', function () {
  result = fooTransform(inputArray, (x: any) => x);
});

When('I filter it with any predicate', function () {
  result = fooFilter(inputArray, () => true);
});

Then('the result should be {string}', function (expectedString: string) {
  // For string results, compare directly. For arrays/objects, parse as JSON
  if (typeof result === 'string') {
    assert.equal(result, expectedString);
  } else {
    const expected = JSON.parse(expectedString);
    assert.deepEqual(result, expected);
  }
});

Then('the original array should remain unchanged', function () {
  assert.deepEqual(inputArray, originalArray, 'Original array was modified');
});

Then('the functions should return new arrays', function () {
  assert.notEqual(result, inputArray, 'Function returned the same array reference');
});

Then('the result should be an empty array', function () {
  assert(Array.isArray(result));
  assert.equal(result.length, 0);
});

// Processing Steps
Given('I have access to the foo processing functions', function () {
  assert(typeof fooProcess === 'function');
  assert(typeof createFooProcessor === 'function');
});

Given('I have a config with prefix {string}', function (prefix: string) {
  config = { prefix };
});

Given('I have a config with prefix {string} and suffix {string}', function (prefix: string, suffix: string) {
  config = { prefix, suffix };
});

When('I create a processor with this config', function () {
  if (!config) throw new Error('Config not set');
  processor = createFooProcessor(config);
});

When('I process the input {string}', function (input: string) {
  if (!config) throw new Error('Config not set');
  // Always create a fresh processor to avoid state bleeding between scenarios
  processor = createFooProcessor(config);
  result = processor.process(input);
});

When('I create processors with this config multiple times', function () {
  if (!config) throw new Error('Config not set');
  processors = [
    createFooProcessor(config),
    createFooProcessor(config),
    createFooProcessor(config)
  ];
});

Then('all processors should behave identically', function () {
  assert(processors.length >= 2, 'Need at least 2 processors to compare');

  const testInput = 'test';
  const results = processors.map(p => p.process(testInput));

  // All results should be identical
  const firstResult = results[0];
  results.forEach((result, index) => {
    assert.equal(result, firstResult, `Processor ${index} returned different result`);
  });
});

Then('processing {string} should always return {string}', function (input: string, expected: string) {
  processors.forEach((proc, index) => {
    const result = proc.process(input);
    assert.equal(result, expected, `Processor ${index} returned unexpected result`);
  });
});

// Greeting Steps
Given('I have access to the foo greeting functions', function () {
  // Verify that greeting functions are available via imports
  assert(typeof fooGreet === 'function');
  assert(typeof createFooProcessor === 'function');
});

Given('I have a greeting config with prefix {string}', function (prefix: string) {
  config = { prefix };
  greetFunction = fooGreet(config);
});

When('I greet {string}', function (name: string) {
  if (!greetFunction) {
    throw new Error('Greeting function not initialized. Set config first.');
  }
  result = greetFunction(name);
});

Then('the result should contain the shared greeting format', function () {
  // Check that result contains "Hello," which is the shared greeting format
  assert(typeof result === 'string', 'Result should be a string');
  assert(result.includes('Hello,'), `Result "${result}" should contain shared greeting format "Hello,"`);
});

Then('the result should start with {string}', function (expectedStart: string) {
  assert(typeof result === 'string', 'Result should be a string');
  assert(result.startsWith(expectedStart), `Result "${result}" should start with "${expectedStart}"`);
});
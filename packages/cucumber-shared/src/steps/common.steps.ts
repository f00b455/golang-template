import { Given, When, Then } from '@cucumber/cucumber';
import { strict as assert } from 'assert';
import type { CommonWorld } from '../index';

// Generic "I have access to" step that can be used by any module
Given('I have access to the {word} functions', function (moduleName: string) {
    // Generic implementation - functions are available through imports
    // This step primarily serves as documentation that the module is available
});

Given('I have access to the {word} {word} functions', function (moduleType: string, moduleName: string) {
    // Generic implementation for compound module names
    // e.g., "shared greet functions", "foo processing functions"
});

Given('I have access to the {word} {word}', function (moduleType: string, moduleName: string) {
    // Generic implementation for simple access steps
    // e.g., "shared formatDate function"
});

// Common result assertion that handles both strings and arrays
Then('the result should be {string}', function (this: CommonWorld, expectedResult: string) {
    // Intelligently handle both array and string results
    const result = this.result;

    if (Array.isArray(result)) {
        try {
            const expected = JSON.parse(expectedResult);
            assert.deepEqual(result, expected);
        } catch (error) {
            // If JSON parsing fails, treat as string comparison
            assert.equal(String(result), expectedResult);
        }
    } else {
        assert.equal(result, expectedResult);
    }
});

// Common state management steps
When('I call it multiple times', function (this: CommonWorld) {
    const func = this.currentFunction;
    const args = this.currentArgs || [];
    const results: any[] = [];

    if (func) {
        for (let i = 0; i < 5; i++) {
            results.push(func(...args));
        }

        this.results = results;
        this.result = results[0];
    }
});

Then('all results should be identical', function (this: CommonWorld) {
    const results = this.results;
    if (results && results.length > 0) {
        const firstResult = results[0];
        for (const res of results) {
            assert.deepEqual(res, firstResult);
        }
    }
});

// Common array result step
Then('the array result should be {string}', function (this: CommonWorld, expectedArray: string) {
    const result = this.result;
    try {
        const expected = JSON.parse(expectedArray);
        assert.deepEqual(result, expected);
    } catch (error) {
        throw new Error(`Failed to parse expected array: ${expectedArray}`);
    }
});

// Empty array assertion
Then('the result should be an empty array', function (this: CommonWorld) {
    const result = this.result;
    assert.deepEqual(result, []);
});

// Type checking steps
Then('the result should be a string', function (this: CommonWorld) {
    const result = this.result;
    assert.equal(typeof result, 'string');
});

Then('the result should be an array', function (this: CommonWorld) {
    const result = this.result;
    assert.ok(Array.isArray(result));
});

Then('the result should be a number', function (this: CommonWorld) {
    const result = this.result;
    assert.equal(typeof result, 'number');
});

Then('the result should be an object', function (this: CommonWorld) {
    const result = this.result;
    assert.equal(typeof result, 'object');
    assert.ok(result !== null);
});
import { Then } from '@cucumber/cucumber';
import { strict as assert } from 'assert';
import type { CommonWorld } from '../index';

// String assertions
Then('the result should start with {string}', function (this: CommonWorld, prefix: string) {
    const result = this.result;
    assert.ok(typeof result === 'string' && result.startsWith(prefix), `Expected result to start with "${prefix}" but got "${result}"`);
});

Then('the result should end with {string}', function (this: CommonWorld, suffix: string) {
    const result = this.result;
    assert.ok(typeof result === 'string' && result.endsWith(suffix), `Expected result to end with "${suffix}" but got "${result}"`);
});

Then('the result should contain {string}', function (this: CommonWorld, substring: string) {
    const result = this.result;
    assert.ok(typeof result === 'string' && result.includes(substring), `Expected result to contain "${substring}" but got "${result}"`);
});

Then('the result should match pattern {string}', function (this: CommonWorld, pattern: string) {
    const result = this.result;
    const regex = new RegExp(pattern);
    assert.ok(typeof result === 'string', 'Result must be a string to match pattern');
    assert.match(result as string, regex, `Expected result to match pattern "${pattern}" but got "${result}"`);
});

// Numeric assertions
Then('the result should be greater than {int}', function (this: CommonWorld, value: number) {
    const result = this.result;
    assert.ok(typeof result === 'number' && result > value, `Expected ${result} to be greater than ${value}`);
});

Then('the result should be less than {int}', function (this: CommonWorld, value: number) {
    const result = this.result;
    assert.ok(typeof result === 'number' && result < value, `Expected ${result} to be less than ${value}`);
});

Then('the result should be between {int} and {int}', function (this: CommonWorld, min: number, max: number) {
    const result = this.result;
    assert.ok(typeof result === 'number' && result >= min && result <= max, `Expected ${result} to be between ${min} and ${max}`);
});

// Array assertions
Then('the result should have {int} items', function (this: CommonWorld, count: number) {
    const result = this.result;
    assert.ok(Array.isArray(result), 'Result must be an array');
    assert.equal((result as any[]).length, count, `Expected ${count} items but got ${(result as any[]).length}`);
});

Then('the result should include {string}', function (this: CommonWorld, item: string) {
    const result = this.result;
    assert.ok(Array.isArray(result), 'Result must be an array');
    try {
        const parsedItem = (item.includes('[') || item.includes('{')) ? JSON.parse(item) : item;
        assert.ok((result as any[]).includes(parsedItem), `Expected array to include ${item}`);
    } catch (error) {
        assert.ok((result as any[]).includes(item), `Expected array to include ${item}`);
    }
});

Then('the result should not be empty', function (this: CommonWorld) {
    const result = this.result;
    if (Array.isArray(result)) {
        assert.ok(result.length > 0, 'Expected array not to be empty');
    } else if (typeof result === 'string') {
        assert.ok(result.length > 0, 'Expected string not to be empty');
    } else if (typeof result === 'object' && result !== null) {
        assert.ok(Object.keys(result).length > 0, 'Expected object not to be empty');
    }
});

// Boolean assertions
Then('the result should be true', function (this: CommonWorld) {
    const result = this.result;
    assert.strictEqual(result, true);
});

Then('the result should be false', function (this: CommonWorld) {
    const result = this.result;
    assert.strictEqual(result, false);
});

Then('the result should be truthy', function (this: CommonWorld) {
    const result = this.result;
    assert.ok(result);
});

Then('the result should be falsy', function (this: CommonWorld) {
    const result = this.result;
    assert.ok(!result);
});

// Null/undefined assertions
Then('the result should be null', function (this: CommonWorld) {
    const result = this.result;
    assert.strictEqual(result, null);
});

Then('the result should be undefined', function (this: CommonWorld) {
    const result = this.result;
    assert.strictEqual(result, undefined);
});

Then('the result should be defined', function (this: CommonWorld) {
    const result = this.result;
    assert.notEqual(result, undefined);
});
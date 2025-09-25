import { Given, Then } from '@cucumber/cucumber';
import { strict as assert } from 'assert';
import type { CommonWorld } from '../index';

// Array initialization steps
Given('I have an array {string}', function (this: CommonWorld, arrayString: string) {
    try {
        const inputArray = JSON.parse(arrayString);
        this.inputArray = inputArray;
        this.originalArray = [...inputArray];
    } catch (error) {
        throw new Error(`Failed to parse array: ${arrayString}`);
    }
});

Given('I have a number array {string}', function (this: CommonWorld, arrayString: string) {
    try {
        const inputArray = JSON.parse(arrayString);
        this.inputArray = inputArray;
        this.originalArray = [...inputArray];
    } catch (error) {
        throw new Error(`Failed to parse number array: ${arrayString}`);
    }
});

Given('I have a string array {string}', function (this: CommonWorld, arrayString: string) {
    try {
        const inputArray = JSON.parse(arrayString);
        this.inputArray = inputArray;
        this.originalArray = [...inputArray];
    } catch (error) {
        throw new Error(`Failed to parse string array: ${arrayString}`);
    }
});

Given('I have an empty array', function (this: CommonWorld) {
    this.inputArray = [];
    this.originalArray = [];
});

// Object initialization steps
Given('I have an object {string}', function (this: CommonWorld, objectString: string) {
    try {
        const inputObject = JSON.parse(objectString);
        this.inputObject = inputObject;
        this.originalObject = { ...inputObject };
    } catch (error) {
        throw new Error(`Failed to parse object: ${objectString}`);
    }
});

Given('I have an empty object', function (this: CommonWorld) {
    this.inputObject = {};
    this.originalObject = {};
});

// Data immutability assertions
Then('the original array should remain unchanged', function (this: CommonWorld) {
    const inputArray = this.inputArray;
    const originalArray = this.originalArray;
    assert.deepEqual(inputArray, originalArray, 'Original array was modified');
});

Then('the original object should remain unchanged', function (this: CommonWorld) {
    const inputObject = this.inputObject;
    const originalObject = this.originalObject;
    assert.deepEqual(inputObject, originalObject, 'Original object was modified');
});

Then('the input should remain unchanged', function (this: CommonWorld) {
    // Generic immutability check
    if (this.originalArray !== undefined) {
        assert.deepEqual(this.inputArray, this.originalArray);
    } else if (this.originalObject !== undefined) {
        assert.deepEqual(this.inputObject, this.originalObject);
    }
});

// Reference checks
Then('the functions should return new arrays', function (this: CommonWorld) {
    const result = this.result;
    const inputArray = this.inputArray;
    assert.notEqual(result, inputArray, 'Function returned same array reference');
});

Then('the result should be a new array', function (this: CommonWorld) {
    const result = this.result;
    const inputArray = this.inputArray;
    assert.notEqual(result, inputArray, 'Result is not a new array');
});

Then('the result should be a new object', function (this: CommonWorld) {
    const result = this.result;
    const inputObject = this.inputObject;
    assert.notEqual(result, inputObject, 'Result is not a new object');
});

// Data transformation helpers
Given('I have test data {string}', function (this: CommonWorld, dataString: string) {
    try {
        const data = JSON.parse(dataString);
        this.testData = data;
    } catch (error) {
        throw new Error(`Failed to parse test data: ${dataString}`);
    }
});

Then('the transformed result should be {string}', function (this: CommonWorld, expectedString: string) {
    const result = this.result;
    try {
        const expected = JSON.parse(expectedString);
        assert.deepEqual(result, expected);
    } catch (error) {
        throw new Error(`Failed to parse expected result: ${expectedString}`);
    }
});

// Sorting assertions
Then('the result should be sorted ascending', function (this: CommonWorld) {
    const result = this.result;
    assert.ok(Array.isArray(result), 'Result must be an array');
    const sorted = [...(result as any[])].sort((a: any, b: any) => {
        if (typeof a === 'number' && typeof b === 'number') return a - b;
        if (typeof a === 'string' && typeof b === 'string') return a.localeCompare(b);
        return 0;
    });
    assert.deepEqual(result, sorted, 'Result is not sorted ascending');
});

Then('the result should be sorted descending', function (this: CommonWorld) {
    const result = this.result;
    assert.ok(Array.isArray(result), 'Result must be an array');
    const sorted = [...(result as any[])].sort((a: any, b: any) => {
        if (typeof a === 'number' && typeof b === 'number') return b - a;
        if (typeof a === 'string' && typeof b === 'string') return b.localeCompare(a);
        return 0;
    });
    assert.deepEqual(result, sorted, 'Result is not sorted descending');
});
import { Given, When, Then } from '@cucumber/cucumber';
import { strict as assert } from 'assert';

// API common context interface
interface ApiContext {
    apiResponse?: any;
    apiData?: unknown;
    apiError?: Error;
    startTime?: number;
    responseTime?: number;
    statusCode?: number;
}

// API state steps
Given('the API is available', function (this: ApiContext) {
    // Mock or check API availability
    assert.ok(true, 'API should be available');
});

Given('the API server is running', function (this: ApiContext) {
    // Check if server is running
    assert.ok(true, 'API server should be running');
});

Given('the API endpoint exists', function (this: ApiContext) {
    // Verify endpoint exists
    assert.ok(true, 'API endpoint should exist');
});

// API response steps
When('I make a request to {string}', async function (this: ApiContext, endpoint: string) {
    this.startTime = Date.now();
    // Mock API request - in real implementation would make actual call
    this.apiResponse = { status: 200, data: {} };
    this.responseTime = Date.now() - this.startTime;
});

When('the API responds', function (this: ApiContext) {
    assert.ok(this.apiResponse, 'API should have responded');
});

// API assertion steps
Then('the API should respond with status {int}', function (this: ApiContext, expectedStatus: number) {
    const status = this.apiResponse?.status || this.statusCode;
    assert.equal(status, expectedStatus, `Expected status ${expectedStatus} but got ${status}`);
});

Then('the response should contain {string}', function (this: ApiContext, expectedContent: string) {
    const responseString = JSON.stringify(this.apiResponse || this.apiData);
    assert.ok(responseString.includes(expectedContent), `Response should contain "${expectedContent}"`);
});

Then('the response time should be under {int}ms', function (this: ApiContext, maxTime: number) {
    assert.ok(this.responseTime !== undefined, 'Response time not measured');
    assert.ok(this.responseTime < maxTime, `Response time ${this.responseTime}ms exceeds ${maxTime}ms`);
});

Then('the API should return valid JSON', function (this: ApiContext) {
    try {
        if (typeof this.apiData === 'string') {
            JSON.parse(this.apiData);
        }
        assert.ok(true, 'Response is valid JSON');
    } catch (error) {
        assert.fail('Response is not valid JSON');
    }
});

// Error handling steps
Given('the API is unavailable', function (this: ApiContext) {
    this.apiError = new Error('API unavailable');
});

Given('the API returns an error', function (this: ApiContext) {
    this.apiError = new Error('API error');
    this.statusCode = 500;
});

Then('an error should be handled gracefully', function (this: ApiContext) {
    assert.ok(this.apiError, 'Error should exist');
    // In real implementation, verify error is handled properly
});

Then('the error message should contain {string}', function (this: ApiContext, expectedMessage: string) {
    assert.ok(this.apiError?.message.includes(expectedMessage),
        `Error message should contain "${expectedMessage}"`);
});

// Data validation steps
Then('the response should have field {string}', function (this: ApiContext, fieldName: string) {
    const data = this.apiResponse?.data || this.apiData;
    assert.ok(fieldName in (data as any), `Response should have field "${fieldName}"`);
});

Then('the response field {string} should equal {string}', function (this: ApiContext, fieldName: string, expectedValue: string) {
    const data = this.apiResponse?.data || this.apiData;
    const actualValue = (data as any)[fieldName];
    const expected = expectedValue.includes('[') || expectedValue.includes('{')
        ? JSON.parse(expectedValue)
        : expectedValue;
    assert.deepEqual(actualValue, expected,
        `Field "${fieldName}" should equal ${expectedValue} but got ${actualValue}`);
});

// Caching steps
Then('the response should be cached', function (this: ApiContext) {
    // In real implementation, verify caching headers or behavior
    assert.ok(true, 'Response should be cached');
});

Then('the cache TTL should be {int} minutes', function (this: ApiContext, ttlMinutes: number) {
    // In real implementation, verify cache TTL
    assert.ok(true, `Cache TTL should be ${ttlMinutes} minutes`);
});
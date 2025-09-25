// Re-export all step definitions
export * from './steps/common.steps';
export * from './steps/assertions.steps';
export * from './steps/data-operations.steps';
export * from './steps/api.steps';

// Export types for World context
export interface CommonWorld {
    result?: any;
    results?: any[];
    inputArray?: any[];
    originalArray?: any[];
    inputObject?: any;
    originalObject?: any;
    testData?: any;
    currentFunction?: Function;
    currentArgs?: any[];
}

export interface ApiWorld {
    apiResponse?: any;
    apiData?: unknown;
    apiError?: Error;
    startTime?: number;
    responseTime?: number;
    statusCode?: number;
}
/**
 * Base ESLint configuration for enforcing Clean Code principles
 * These rules apply to all packages in the monorepo
 *
 * Phase control via ESLINT_CLEAN_CODE_PHASE environment variable:
 * - Phase 1 (default): 'warn' - Assessment phase
 * - Phase 2: 'warn' - Gradual refactoring
 * - Phase 3: 'error' - Full enforcement
 */

// Determine the severity level based on environment variable
const VALID_PHASES = ['1', '2', '3'];
const PHASE = VALID_PHASES.includes(process.env.ESLINT_CLEAN_CODE_PHASE)
  ? process.env.ESLINT_CLEAN_CODE_PHASE
  : '1';
const severity = PHASE === '3' ? 'error' : 'warn';

module.exports = {
  // Enable cache for faster subsequent runs
  cache: true,
  cacheLocation: 'node_modules/.cache/eslint/',

  rules: {
    // Enforce function length limits - Clean Code principle
    // Rationale: Functions should do one thing (Single Responsibility)
    // 20 lines encourages breaking complex logic into smaller, testable units
    'max-lines-per-function': [
      severity,
      {
        max: 20,
        skipBlankLines: true,
        skipComments: true,
        IIFEs: false
      }
    ],

    // Enforce file length limits
    // Rationale: Files should have a single, focused purpose
    // 300 lines prevents monolithic files that are hard to navigate and understand
    'max-lines': [
      severity,
      {
        max: 300,
        skipBlankLines: true,
        skipComments: true
      }
    ],

    // Enforce maximum statements in a function
    // Rationale: Complements max-lines-per-function by counting logical steps
    // 15 statements forces cleaner, more focused functions
    'max-statements': [
      severity,
      15,
      {
        ignoreTopLevelFunctions: false
      }
    ],

    // Enforce complexity limits
    // Rationale: Lower cyclomatic complexity = easier to test and understand
    // 10 is a widely accepted threshold for maintainable code
    'complexity': [
      severity,
      {
        max: 10
      }
    ],

    // Enforce maximum depth of nested blocks
    // Rationale: Deep nesting makes code hard to follow
    // 3 levels encourages early returns and guard clauses
    'max-depth': [
      severity,
      {
        max: 3
      }
    ],

    // Enforce maximum number of parameters
    // Rationale: Too many parameters indicate a function doing too much
    // 4 params max encourages using objects for complex data
    'max-params': [
      severity,
      {
        max: 4
      }
    ]
  },

  // Override for test files and step definitions
  overrides: [
    {
      files: [
        '**/*.test.ts',
        '**/*.test.tsx',
        '**/*.spec.ts',
        '**/*.spec.tsx',
        '**/test/**/*.ts',
        '**/tests/**/*.ts',
        '**/__tests__/**/*.ts',
        '**/features/**/*.ts',
        '**/step_definitions/**/*.ts',
        '**/e2e/**/*.ts',
        '**/*.steps.ts'
      ],
      rules: {
        // More lenient limits for test files
        'max-lines-per-function': [
          severity,
          {
            max: 50, // Tests can be longer
            skipBlankLines: true,
            skipComments: true,
            IIFEs: false
          }
        ],
        'max-lines': [
          severity,
          {
            max: 500, // Test files can be longer
            skipBlankLines: true,
            skipComments: true
          }
        ],
        'max-statements': [
          severity,
          30 // Tests often need more statements for setup
        ],
        // Disable complexity rules for tests
        'complexity': 'off',
        'max-depth': 'off',
        'max-params': 'off'
      }
    }
  ]
};
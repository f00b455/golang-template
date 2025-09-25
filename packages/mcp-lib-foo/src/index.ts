import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import { CallToolRequestSchema, ListToolsRequestSchema } from '@modelcontextprotocol/sdk/types.js';
import {
  fooProcess,
  fooGreet,
  createFooProcessor,
  fooTransform,
  fooFilter,
  type FooConfig
} from '@ts-template/lib-foo';

const server = new Server(
  {
    name: 'mcp-lib-foo',
    version: '0.0.0',
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// Tool definitions with JSON schemas
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: 'fooProcess',
        description: 'Process text with prefix and optional suffix',
        inputSchema: {
          type: 'object',
          properties: {
            input: {
              type: 'string',
              description: 'The input text to process',
            },
            prefix: {
              type: 'string',
              description: 'Prefix to add to the input',
            },
            suffix: {
              type: 'string',
              description: 'Optional suffix to add to the input',
            },
          },
          required: ['input', 'prefix'],
        },
      },
      {
        name: 'fooGreet',
        description: 'Create greetings with foo processing',
        inputSchema: {
          type: 'object',
          properties: {
            name: {
              type: 'string',
              description: 'The name to greet',
            },
            prefix: {
              type: 'string',
              description: 'Prefix for the greeting',
            },
            suffix: {
              type: 'string',
              description: 'Optional suffix for the greeting',
            },
          },
          required: ['name', 'prefix'],
        },
      },
      {
        name: 'createFooProcessor',
        description: 'Create a configured foo processor',
        inputSchema: {
          type: 'object',
          properties: {
            prefix: {
              type: 'string',
              description: 'Default prefix for the processor',
            },
            suffix: {
              type: 'string',
              description: 'Optional default suffix for the processor',
            },
            action: {
              type: 'string',
              enum: ['process', 'greetWithFoo'],
              description: 'Action to perform with the processor',
            },
            input: {
              type: 'string',
              description: 'Input for the selected action',
            },
          },
          required: ['prefix', 'action', 'input'],
        },
      },
      {
        name: 'fooTransform',
        description: 'Transform an array of strings with a custom transformer',
        inputSchema: {
          type: 'object',
          properties: {
            data: {
              type: 'array',
              items: {
                type: 'string',
              },
              description: 'Array of strings to transform',
            },
            transformType: {
              type: 'string',
              enum: ['uppercase', 'lowercase', 'reverse', 'trim'],
              description: 'Type of transformation to apply',
            },
          },
          required: ['data', 'transformType'],
        },
      },
      {
        name: 'fooFilter',
        description: 'Filter an array based on a predicate',
        inputSchema: {
          type: 'object',
          properties: {
            items: {
              type: 'array',
              items: {
                oneOf: [
                  { type: 'string' },
                  { type: 'number' },
                  { type: 'boolean' },
                ],
              },
              description: 'Array of items to filter',
            },
            filterType: {
              type: 'string',
              enum: ['truthy', 'falsy', 'nonEmpty', 'numeric', 'string'],
              description: 'Type of filter to apply',
            },
            customPredicate: {
              type: 'string',
              description: 'Optional custom predicate (for strings: minLength:N)',
            },
          },
          required: ['items', 'filterType'],
        },
      },
    ],
  };
});

// Tool handlers
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  try {
    switch (name) {
      case 'fooProcess': {
        const { input, prefix, suffix } = args as { input: string; prefix: string; suffix?: string };
        const config: FooConfig = suffix !== undefined ? { prefix, suffix } : { prefix };
        const result = fooProcess(config)(input);
        return {
          content: [
            {
              type: 'text',
              text: result,
            },
          ],
        };
      }

      case 'fooGreet': {
        const { name: userName, prefix, suffix } = args as { name: string; prefix: string; suffix?: string };
        const config: FooConfig = suffix !== undefined ? { prefix, suffix } : { prefix };
        const result = fooGreet(config)(userName);
        return {
          content: [
            {
              type: 'text',
              text: result,
            },
          ],
        };
      }

      case 'createFooProcessor': {
        const { prefix, suffix, action, input } = args as {
          prefix: string;
          suffix?: string;
          action: 'process' | 'greetWithFoo';
          input: string;
        };
        const config: FooConfig = suffix !== undefined ? { prefix, suffix } : { prefix };
        const processor = createFooProcessor(config);

        let result: string;
        if (action === 'process') {
          result = processor.process(input);
        } else {
          result = processor.greetWithFoo(input);
        }

        return {
          content: [
            {
              type: 'text',
              text: result,
            },
          ],
        };
      }

      case 'fooTransform': {
        const { data, transformType } = args as {
          data: string[];
          transformType: 'uppercase' | 'lowercase' | 'reverse' | 'trim';
        };

        let transformer: (value: string) => string;
        switch (transformType) {
          case 'uppercase':
            transformer = (s: string) => s.toUpperCase();
            break;
          case 'lowercase':
            transformer = (s: string) => s.toLowerCase();
            break;
          case 'reverse':
            transformer = (s: string) => s.split('').reverse().join('');
            break;
          case 'trim':
            transformer = (s: string) => s.trim();
            break;
          default:
            throw new Error(`Unknown transform type: ${transformType}`);
        }

        const result = fooTransform(data, transformer);
        return {
          content: [
            {
              type: 'text',
              text: JSON.stringify(result),
            },
          ],
        };
      }

      case 'fooFilter': {
        const { items, filterType, customPredicate } = args as {
          items: unknown[];
          filterType: 'truthy' | 'falsy' | 'nonEmpty' | 'numeric' | 'string';
          customPredicate?: string;
        };

        let predicate: (value: unknown) => boolean;
        switch (filterType) {
          case 'truthy':
            predicate = (v: unknown) => !!v;
            break;
          case 'falsy':
            predicate = (v: unknown) => !v;
            break;
          case 'nonEmpty':
            predicate = (v: unknown) => {
              if (typeof v === 'string') return v.length > 0;
              if (Array.isArray(v)) return v.length > 0;
              return v != null;
            };
            break;
          case 'numeric':
            predicate = (v: unknown) => typeof v === 'number' && !isNaN(v as number);
            break;
          case 'string':
            if (customPredicate && customPredicate.startsWith('minLength:')) {
              const parts = customPredicate.split(':');
              const minLength = parts[1] ? parseInt(parts[1], 10) : 0;
              predicate = (v: unknown) => typeof v === 'string' && v.length >= minLength;
            } else {
              predicate = (v: unknown) => typeof v === 'string';
            }
            break;
          default:
            throw new Error(`Unknown filter type: ${filterType}`);
        }

        const result = fooFilter(items, predicate);
        return {
          content: [
            {
              type: 'text',
              text: JSON.stringify(result),
            },
          ],
        };
      }

      default:
        throw new Error(`Unknown tool: ${name}`);
    }
  } catch (error) {
    return {
      content: [
        {
          type: 'text',
          text: `Error: ${error instanceof Error ? error.message : String(error)}`,
        },
      ],
    };
  }
});

// Start the server
async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error('MCP Server for lib-foo started');
}

main().catch((error) => {
  console.error('Failed to start server:', error);
  process.exit(1);
});
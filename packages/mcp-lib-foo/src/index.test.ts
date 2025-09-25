import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

// Mock all dependencies BEFORE any imports
vi.mock('@ts-template/lib-foo', () => ({
  fooProcess: vi.fn(),
  fooGreet: vi.fn(),
  createFooProcessor: vi.fn(),
  fooTransform: vi.fn(),
  fooFilter: vi.fn(),
}));

vi.mock('@modelcontextprotocol/sdk/server/index.js', () => ({
  Server: vi.fn(),
}));

vi.mock('@modelcontextprotocol/sdk/server/stdio.js', () => ({
  StdioServerTransport: vi.fn().mockImplementation(() => ({
    connect: vi.fn(),
  })),
}));

// Mock process.exit to prevent test termination
const mockExit = vi.fn();
vi.stubGlobal('process', { exit: mockExit });

describe('MCP Server for lib-foo', () => {
  let mockServer: { setRequestHandler: ReturnType<typeof vi.fn>; connect: ReturnType<typeof vi.fn> };
  let mockSetRequestHandler: ReturnType<typeof vi.fn>;
  let mockConnect: ReturnType<typeof vi.fn>;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let Server: any;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let StdioServerTransport: any;

  beforeEach(async () => {
    // Clear module cache to ensure fresh imports
    vi.resetModules();

    // Reset all mocks
    vi.clearAllMocks();

    // Set up fresh mocks
    mockSetRequestHandler = vi.fn();
    mockConnect = vi.fn();
    mockServer = {
      setRequestHandler: mockSetRequestHandler,
      connect: mockConnect,
    };

    // Import mocked modules
    const serverModule = await import('@modelcontextprotocol/sdk/server/index.js');
    const transportModule = await import('@modelcontextprotocol/sdk/server/stdio.js');
    Server = serverModule.Server;
    StdioServerTransport = transportModule.StdioServerTransport;

    // Configure mocks
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (Server as any).mockImplementation(() => mockServer);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (StdioServerTransport as any).mockImplementation(() => ({
      connect: mockConnect,
    }));
  });

  afterEach(() => {
    vi.resetModules();
  });

  describe('Server Initialization', () => {
    it('should create server with correct configuration', async () => {
      // Import the module to trigger server creation
      await import('./index.js');

      expect(Server).toHaveBeenCalledWith(
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
    });

    it('should register handlers and connect', async () => {
      await import('./index.js');

      // Should register both handlers
      expect(mockSetRequestHandler).toHaveBeenCalledTimes(2);

      // Should create transport and connect
      expect(StdioServerTransport).toHaveBeenCalled();
      expect(mockConnect).toHaveBeenCalled();
    });
  });

  describe('Tools Registration', () => {
    it('should register ListToolsRequestSchema handler', async () => {
      const { ListToolsRequestSchema } = await import('@modelcontextprotocol/sdk/types.js');
      await import('./index.js');

      const listToolsCalls = mockSetRequestHandler.mock.calls.filter(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (call: any) => call[0] === ListToolsRequestSchema
      );
      expect(listToolsCalls).toHaveLength(1);
    });

    it('should register CallToolRequestSchema handler', async () => {
      const { CallToolRequestSchema } = await import('@modelcontextprotocol/sdk/types.js');
      await import('./index.js');

      const callToolCalls = mockSetRequestHandler.mock.calls.filter(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (call: any) => call[0] === CallToolRequestSchema
      );
      expect(callToolCalls).toHaveLength(1);
    });

    it('should expose 5 tools when listing', async () => {
      const { ListToolsRequestSchema } = await import('@modelcontextprotocol/sdk/types.js');
      await import('./index.js');

      // Find the tools list handler
      const listToolsCall = mockSetRequestHandler.mock.calls.find(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (call: any) => call[0] === ListToolsRequestSchema
      );
      expect(listToolsCall).toBeDefined();
      if (!listToolsCall) throw new Error('Tools list handler not found');

      const toolsHandler = listToolsCall[1];
      const result = await toolsHandler();

      expect(result.tools).toHaveLength(5);

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const toolNames = result.tools.map((t: any) => t.name);
      expect(toolNames).toContain('fooProcess');
      expect(toolNames).toContain('fooGreet');
      expect(toolNames).toContain('createFooProcessor');
      expect(toolNames).toContain('fooTransform');
      expect(toolNames).toContain('fooFilter');
    });
  });

  describe('Tool Execution', () => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    let callHandler: any;

    beforeEach(async () => {
      const libFoo = await import('@ts-template/lib-foo');
      const { CallToolRequestSchema } = await import('@modelcontextprotocol/sdk/types.js');

      // Set up mock return values
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (libFoo.fooProcess as any).mockReturnValue(() => 'mocked-process-result');
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (libFoo.fooGreet as any).mockReturnValue(() => 'mocked-greet-result');
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (libFoo.createFooProcessor as any).mockReturnValue({
        process: () => 'mocked-processor-result',
        greetWithFoo: () => 'mocked-greeting-result',
      });
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (libFoo.fooTransform as any).mockReturnValue(['MOCKED', 'RESULT']);
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (libFoo.fooFilter as any).mockReturnValue(['filtered', 'result']);

      await import('./index.js');

      // Find the call tool handler
      const callHandlerCall = mockSetRequestHandler.mock.calls.find(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (call: any) => call[0] === CallToolRequestSchema
      );
      expect(callHandlerCall).toBeDefined();
      if (!callHandlerCall) throw new Error('Call tool handler not found');
      callHandler = callHandlerCall[1];
    });

    it('should execute fooProcess tool correctly', async () => {
      expect(callHandler).toBeTruthy();

      const libFoo = await import('@ts-template/lib-foo');

      const request = {
        params: {
          name: 'fooProcess',
          arguments: {
            input: 'test-input',
            prefix: 'pre-',
            suffix: '-post',
          },
        },
      };

      const result = await callHandler(request);

      expect(libFoo.fooProcess).toHaveBeenCalledWith({ prefix: 'pre-', suffix: '-post' });
      expect(result.content).toHaveLength(1);
      expect(result.content[0].type).toBe('text');
      expect(result.content[0].text).toBe('mocked-process-result');
    });

    it('should handle optional suffix parameter', async () => {
      const libFoo = await import('@ts-template/lib-foo');

      const request = {
        params: {
          name: 'fooProcess',
          arguments: {
            input: 'test-input',
            prefix: 'pre-',
          },
        },
      };

      await callHandler(request);

      expect(libFoo.fooProcess).toHaveBeenCalledWith({ prefix: 'pre-' });
    });

    it('should execute fooGreet tool correctly', async () => {
      const libFoo = await import('@ts-template/lib-foo');

      const request = {
        params: {
          name: 'fooGreet',
          arguments: {
            name: 'Alice',
            prefix: 'Hello-',
            suffix: '!',
          },
        },
      };

      const result = await callHandler(request);

      expect(libFoo.fooGreet).toHaveBeenCalledWith({ prefix: 'Hello-', suffix: '!' });
      expect(result.content[0].text).toBe('mocked-greet-result');
    });

    it('should execute fooTransform tool correctly', async () => {
      const libFoo = await import('@ts-template/lib-foo');

      const request = {
        params: {
          name: 'fooTransform',
          arguments: {
            data: ['hello', 'world'],
            transformType: 'uppercase',
          },
        },
      };

      const result = await callHandler(request);

      // fooTransform should be called with data and a transformer function
      expect(libFoo.fooTransform).toHaveBeenCalledWith(['hello', 'world'], expect.any(Function));
      expect(result.content[0].text).toBe('["MOCKED","RESULT"]');
    });

    it('should handle invalid transform type', async () => {
      const request = {
        params: {
          name: 'fooTransform',
          arguments: {
            data: ['test'],
            transformType: 'invalid',
          },
        },
      };

      const result = await callHandler(request);

      expect(result.content[0].text).toContain('Error:');
    });

    it('should execute fooFilter tool correctly', async () => {
      const libFoo = await import('@ts-template/lib-foo');

      const request = {
        params: {
          name: 'fooFilter',
          arguments: {
            items: [1, 'test', 0, '', 'valid', null],
            filterType: 'truthy',
          },
        },
      };

      const result = await callHandler(request);

      // fooFilter should be called with items and a predicate function
      expect(libFoo.fooFilter).toHaveBeenCalledWith(
        [1, 'test', 0, '', 'valid', null],
        expect.any(Function)
      );
      expect(result.content[0].text).toBe('["filtered","result"]');
    });

    it('should handle unknown tool name', async () => {
      const request = {
        params: {
          name: 'unknownTool',
          arguments: {},
        },
      };

      const result = await callHandler(request);

      expect(result.content[0].text).toContain('Unknown tool: unknownTool');
    });

    it('should handle errors from lib-foo functions', async () => {
      const libFoo = await import('@ts-template/lib-foo');
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (libFoo.fooProcess as any).mockImplementation(() => {
        throw new Error('Mock error from lib-foo');
      });

      const request = {
        params: {
          name: 'fooProcess',
          arguments: {
            input: 'test',
            prefix: 'pre-',
          },
        },
      };

      const result = await callHandler(request);

      expect(result.content[0].text).toContain('Error: Mock error from lib-foo');
    });
  });
});
# Multi-stage Dockerfile for MCP Server
FROM node:20-alpine AS base
RUN apk add --no-cache libc6-compat
RUN corepack enable && corepack prepare pnpm@latest --activate
WORKDIR /app

# Dependencies stage
FROM base AS deps
COPY pnpm-lock.yaml ./
COPY pnpm-workspace.yaml ./
COPY package.json ./
COPY packages/mcp-lib-foo/package.json ./packages/mcp-lib-foo/
COPY packages/lib-foo/package.json ./packages/lib-foo/
COPY packages/shared/package.json ./packages/shared/
RUN pnpm install --frozen-lockfile --filter @ts-template/mcp-lib-foo...

# Build stage
FROM base AS build
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/packages ./packages
COPY tsconfig.base.json ./
COPY packages/shared ./packages/shared
COPY packages/lib-foo ./packages/lib-foo
COPY packages/mcp-lib-foo ./packages/mcp-lib-foo

# Build dependencies first
WORKDIR /app/packages/shared
RUN pnpm build

WORKDIR /app/packages/lib-foo
RUN pnpm build

WORKDIR /app/packages/mcp-lib-foo
RUN pnpm build

# Production stage
FROM node:20-alpine AS runner
RUN apk add --no-cache libc6-compat
RUN corepack enable && corepack prepare pnpm@latest --activate
WORKDIR /app

# Create non-root user
RUN addgroup --system --gid 1001 mcp
RUN adduser --system --uid 1001 mcp

# Copy necessary files
COPY --from=build --chown=mcp:mcp /app/packages/mcp-lib-foo/dist ./dist
COPY --from=build --chown=mcp:mcp /app/packages/mcp-lib-foo/package.json ./package.json
COPY --from=build --chown=mcp:mcp /app/node_modules ./node_modules

USER mcp

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "process.exit(0)" || exit 1

# Environment variables
ENV NODE_ENV=production
ENV MCP_SERVER_NAME=mcp-lib-foo
ENV MCP_SERVER_VERSION=0.0.0

# Start server
EXPOSE 3000
CMD ["node", "dist/index.js"]
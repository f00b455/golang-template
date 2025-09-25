import { FastifyPluginAsync } from 'fastify';
import { greet } from '@ts-template/shared';

interface GreetQuery {
  name?: string;
}

export const greetRoute: FastifyPluginAsync = async function (fastify) {
  fastify.get<{ Querystring: GreetQuery }>(
    '/greet',
    {
      schema: {
        description: 'Greet endpoint',
        tags: ['greet'],
        querystring: {
          type: 'object',
          properties: {
            name: {
              type: 'string',
              description: 'Name to greet',
            },
          },
        },
        response: {
          200: {
            type: 'object',
            properties: {
              message: {
                type: 'string',
              },
            },
          },
        },
      },
    },
    async function (request) {
      const { name = 'World' } = request.query;
      return { message: greet(name) };
    }
  );
};
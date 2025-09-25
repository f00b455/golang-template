import { greet } from '@ts-template/shared';

export interface FooConfig {
  readonly prefix: string;
  readonly suffix?: string;
}

// Pure functional approach instead of class
export const fooProcess = (config: Readonly<FooConfig>) => (input: string): string => {
  return `${config.prefix}${input}${config.suffix || ''}`;
};

export const fooGreet = (config: Readonly<FooConfig>) => (name: string): string => {
  const greeting = greet(name);
  return fooProcess(config)(greeting);
};

export const createFooProcessor = (config: Readonly<FooConfig>) => {
  return {
    process: fooProcess(config),
    greetWithFoo: fooGreet(config),
  };
};

export const fooTransform = (
  data: readonly string[],
  transformer: (value: string) => string
): readonly string[] => {
  return data.map(transformer);
};

export const fooFilter = <T>(
  items: readonly T[],
  predicate: (value: T) => boolean
): readonly T[] => {
  return items.filter(predicate);
};
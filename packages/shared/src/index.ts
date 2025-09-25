export const greet = (name: string): string => {
  if (!name.trim()) {
    return 'Error: Name cannot be empty';
  }
  return `Hello, ${name}!`;
};

export const formatDate = (date: Readonly<Date>): string => {
  const parts = date.toISOString().split('T');
  return parts[0] ?? '';
};

export type User = Readonly<{
  readonly id: string;
  readonly name: string;
  readonly email: string;
  readonly createdAt: Date;
}>;

export type RssHeadline = Readonly<{
  readonly title: string;
  readonly link: string;
  readonly publishedAt: string;
  readonly source: string;
}>;

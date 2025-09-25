import type { Metadata } from 'next';
import { GreetForm } from '@/components/GreetForm';
import { NewsHeadlines } from '@/components/NewsHeadlines';

export const metadata: Metadata = {
  title: 'TypeScript Template',
  description: 'A clean TypeScript monorepo template',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}): JSX.Element {
  return (
    <html lang="en">
      <body className="min-h-screen bg-gray-50">
        <header className="bg-white shadow-sm border-b">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
            <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between space-y-4 lg:space-y-0 lg:space-x-8">
              <h1 className="text-3xl font-bold text-gray-900">
                TypeScript Monorepo Template
              </h1>
              <div className="lg:max-w-md xl:max-w-lg">
                <NewsHeadlines />
              </div>
            </div>
          </div>
        </header>

        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <div>
              <GreetForm />
            </div>
            <div>{children}</div>
          </div>
        </main>
      </body>
    </html>
  );
}

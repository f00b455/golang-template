import { greet } from '@ts-template/shared';

export default function Home(): JSX.Element {
  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold mb-4 text-gray-800">
        Static Greeting foo
      </h2>
      <div className="p-4 bg-blue-50 border border-blue-200 rounded-md">
        <p className="text-blue-800 font-medium text-lg">
          {greet('TypeScript Monorepo')}
        </p>
      </div>
      <div className="mt-4 text-sm text-gray-600">
        <p>This greeting is generated server-side using the shared package.</p>
        <p>Compare it with the dynamic API greeting on the left!</p>
      </div>
    </div>
  );
}

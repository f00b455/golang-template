const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3001';

interface GreetResponse {
  message: string;
}

export const apiClient = {
  greet: async (name?: string): Promise<GreetResponse> => {
    const url = new URL('/api/greet', API_BASE_URL);
    if (name) {
      url.searchParams.set('name', name);
    }

    const response = await fetch(url.toString());

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  },
};
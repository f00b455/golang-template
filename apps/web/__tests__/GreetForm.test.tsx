import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { GreetForm } from '@/components/GreetForm';

// Mock the API client
vi.mock('@/lib/api', () => ({
  apiClient: {
    greet: vi.fn(),
  },
}));

import { apiClient } from '@/lib/api';
const mockApiClient = vi.mocked(apiClient);

describe('GreetForm', () => {

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render form elements', () => {
    render(<GreetForm />);

    expect(screen.getByRole('heading', { name: 'API Greeting' })).toBeInTheDocument();
    expect(screen.getByLabelText(/name \(optional\)/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Greet' })).toBeInTheDocument();
  });

  it('should call API with default when form is submitted without name', async () => {
    mockApiClient.greet.mockResolvedValue({ message: 'Hello, World!' });

    render(<GreetForm />);

    const submitButton = screen.getByRole('button', { name: 'Greet' });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockApiClient.greet).toHaveBeenCalledWith(undefined);
    });

    expect(screen.getByText('Hello, World!')).toBeInTheDocument();
  });

  it('should call API with name when form is submitted with name', async () => {
    mockApiClient.greet.mockResolvedValue({ message: 'Hello, Alice!' });

    render(<GreetForm />);

    const nameInput = screen.getByLabelText(/name \(optional\)/i);
    fireEvent.change(nameInput, { target: { value: 'Alice' } });

    const submitButton = screen.getByRole('button', { name: 'Greet' });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockApiClient.greet).toHaveBeenCalledWith('Alice');
    });

    expect(screen.getByText('Hello, Alice!')).toBeInTheDocument();
  });

  it('should show loading state during API call', async () => {
    mockApiClient.greet.mockImplementation(() => new Promise(resolve =>
      setTimeout(() => resolve({ message: 'Hello!' }), 100)
    ));

    render(<GreetForm />);

    const submitButton = screen.getByRole('button', { name: 'Greet' });
    fireEvent.click(submitButton);

    expect(screen.getByRole('button', { name: 'Loading...' })).toBeDisabled();

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Greet' })).toBeEnabled();
    });
  });

  it('should show error message when API call fails', async () => {
    mockApiClient.greet.mockRejectedValue(new Error('Network error'));

    render(<GreetForm />);

    const submitButton = screen.getByRole('button', { name: 'Greet' });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText('Network error')).toBeInTheDocument();
    });
  });

  it('should clear error when new successful request is made', async () => {
    mockApiClient.greet
      .mockRejectedValueOnce(new Error('Network error'))
      .mockResolvedValueOnce({ message: 'Hello, World!' });

    render(<GreetForm />);

    const submitButton = screen.getByRole('button', { name: 'Greet' });

    // First call fails
    fireEvent.click(submitButton);
    await waitFor(() => {
      expect(screen.getByText('Network error')).toBeInTheDocument();
    });

    // Second call succeeds
    fireEvent.click(submitButton);
    await waitFor(() => {
      expect(screen.queryByText('Network error')).not.toBeInTheDocument();
      expect(screen.getByText('Hello, World!')).toBeInTheDocument();
    });
  });
});
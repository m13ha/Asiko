import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useAppointmentByAppCode } from './hooks';
import * as api from './api';

// Mock the API module
jest.mock('./api', () => ({
  getAppointmentByAppCode: jest.fn(),
}));

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      cacheTime: 0,
    },
  },
});

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
);

describe('useAppointmentByAppCode', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    queryClient.clear();
  });

  it('should fetch appointment by app code successfully', async () => {
    const mockAppointment = {
      id: 'test-id',
      appCode: 'TEST123',
      title: 'Test Appointment',
      antiScalpingLevel: 'strict' as const,
      type: 'single' as const,
      status: 'pending' as const,
      bookingDuration: 30,
      startDate: '2024-01-01T00:00:00Z',
      endDate: '2024-01-01T00:00:00Z',
      startTime: '2024-01-01T10:00:00Z',
      endTime: '2024-01-01T11:00:00Z',
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    };

    (api.getAppointmentByAppCode as jest.Mock).mockResolvedValue(mockAppointment);

    const { result } = renderHook(() => useAppointmentByAppCode('TEST123'), {
      wrapper,
    });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data).toEqual(mockAppointment);
    expect(api.getAppointmentByAppCode).toHaveBeenCalledWith('TEST123');
  });

  it('should not fetch if appCode is empty', async () => {
    const { result } = renderHook(() => useAppointmentByAppCode(''), {
      wrapper,
    });

    // The query should be disabled when appCode is empty
    expect(result.current.isFetching).toBe(false);
    expect(result.current.isFetched).toBe(false);
  });

  it('should handle API errors', async () => {
    const mockError = new Error('Failed to fetch appointment');
    (api.getAppointmentByAppCode as jest.Mock).mockRejectedValue(mockError);

    const { result } = renderHook(() => useAppointmentByAppCode('TEST123'), {
      wrapper,
    });

    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });

    expect(result.current.error).toEqual(mockError);
  });
});
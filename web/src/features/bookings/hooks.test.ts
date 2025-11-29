import { renderHook, waitFor, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useDeviceToken } from '../auth/hooks';
import { useAppointmentByAppCode } from '../appointments/hooks';
import { useBookGuest } from './hooks';
import * as bookingApi from './api';
import * as appointmentApi from '../appointments/api';
import * as authApi from '../auth/api';

// Mock the API modules
jest.mock('./api', () => ({
  bookGuest: jest.fn(),
}));

jest.mock('../appointments/api', () => ({
  getAppointmentByAppCode: jest.fn(),
}));

jest.mock('../auth/api', () => ({
  authApi: {
    generateDeviceToken: jest.fn(),
  },
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

describe('Booking flow with device token', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    queryClient.clear();
  });

  it('should generate device token for strict anti-scalping appointments', async () => {
    const mockDeviceTokenResponse = { device_token: 'test-device-token-123' };
    const mockAppointment = {
      id: 'test-id',
      appCode: 'STRICT123',
      title: 'Strict Appointment',
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

    (authApi.authApi.generateDeviceToken as jest.Mock).mockResolvedValue(mockDeviceTokenResponse);
    (appointmentApi.getAppointmentByAppCode as jest.Mock).mockResolvedValue(mockAppointment);

    // Testing the device token hook
    const { result } = renderHook(() => useDeviceToken(), { wrapper });

    // Execute the mutation
    act(() => {
      result.current.mutate({ deviceId: 'test-device-id' });
    });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(authApi.authApi.generateDeviceToken).toHaveBeenCalledWith({
      device: { deviceId: 'test-device-id' },
    });
    expect(result.current.data).toEqual(mockDeviceTokenResponse);
  });

  it('should make booking request with device token when required', async () => {
    const mockBookingResponse = {
      id: 'booking-id',
      appCode: 'STRICT123',
      bookingCode: 'BOOK123',
      name: 'Test User',
      email: 'test@example.com',
      status: 'active',
      date: '2024-01-01T00:00:00Z',
      startTime: '2024-01-01T10:00:00Z',
      endTime: '2024-01-01T11:00:00Z',
    };

    (bookingApi.bookGuest as jest.Mock).mockResolvedValue(mockBookingResponse);

    const { result } = renderHook(() => useBookGuest(), { wrapper });

    const bookingPayload = {
      appCode: 'STRICT123',
      date: '2024-01-01T00:00:00Z',
      startTime: '2024-01-01T10:00:00Z',
      endTime: '2024-01-01T11:00:00Z',
      attendeeCount: 1,
      name: 'Test User',
      email: 'test@example.com',
      deviceToken: 'test-device-token-123',
    };

    act(() => {
      result.current.mutate(bookingPayload);
    });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(bookingApi.bookGuest).toHaveBeenCalledWith(bookingPayload);
    expect(result.current.data).toEqual(mockBookingResponse);
  });

  it('should handle device token generation error gracefully', async () => {
    const mockError = new Error('Failed to generate device token');
    (authApi.authApi.generateDeviceToken as jest.Mock).mockRejectedValue(mockError);

    const { result } = renderHook(() => useDeviceToken(), { wrapper });

    act(() => {
      result.current.mutate({ deviceId: 'test-device-id' });
    });

    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });

    expect(result.current.error).toEqual(mockError);
  });
});
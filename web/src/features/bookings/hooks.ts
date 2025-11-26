import { useMutation, useQuery } from '@tanstack/react-query';
import * as API from '@appointment-master/api-client';
import type { RequestsBookingRequest } from '@appointment-master/api-client';
import toast from 'react-hot-toast';
import { bookGuest, bookRegistered, cancelBookingByCode, getAvailableSlots, getAvailableSlotsByDay, getBookingByCode, getMyRegisteredBookings, updateBookingByCode, rejectBookingByCode } from './api';
import { useQueryClient } from '@tanstack/react-query';

async function parseError(e: unknown): Promise<string> {
  if (e instanceof API.ResponseError) {
    try { const body = await e.response.json(); return body?.message || e.response.statusText || 'Request failed'; }
    catch { return e.response.statusText || 'Request failed'; }
  }
  return 'Something went wrong';
}

export function useAvailableSlots(appCode: string, params?: { page?: number; size?: number }) {
  const page = params?.page ?? 1;
  const size = params?.size ?? 500;
  return useQuery({ 
    queryKey: ['slots', appCode, page, size], 
    queryFn: () => getAvailableSlots(appCode, { page, size }), 
    enabled: !!appCode 
  });
}

export function useAvailableSlotsByDay(appCode: string, date: string, params?: { page?: number; size?: number }) {
  const page = params?.page ?? 1;
  const size = params?.size ?? 200;
  return useQuery({ 
    queryKey: ['slots-by-day', appCode, date, page, size], 
    queryFn: () => getAvailableSlotsByDay(appCode, date, { page, size }), 
    enabled: !!appCode && !!date 
  });
}

export function useBookGuest() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: RequestsBookingRequest) => bookGuest(payload),
    onSuccess: async () => {
      toast.success('Booking confirmed');
      await Promise.all([
        qc.invalidateQueries({ queryKey: ['slots'] }),
        qc.invalidateQueries({ queryKey: ['slots-by-day'] }),
        qc.invalidateQueries({ queryKey: ['my-bookings'] }),
      ]);
    },
    onError: async (e) => toast.error(await parseError(e)),
  });
}

export function useBookRegistered() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: RequestsBookingRequest) => bookRegistered(payload),
    onSuccess: async () => {
      toast.success('Booking confirmed');
      await Promise.all([
        qc.invalidateQueries({ queryKey: ['slots'] }),
        qc.invalidateQueries({ queryKey: ['slots-by-day'] }),
        qc.invalidateQueries({ queryKey: ['my-bookings'] }),
      ]);
    },
    onError: async (e) => toast.error(await parseError(e)),
  });
}

export function useBookingByCode(bookingCode: string) {
  return useQuery({ queryKey: ['booking', bookingCode], queryFn: () => getBookingByCode(bookingCode), enabled: !!bookingCode });
}

export function useUpdateBooking(bookingCode: string) {
  return useMutation({
    mutationFn: (payload: RequestsBookingRequest) => updateBookingByCode(bookingCode, payload),
    onSuccess: () => toast.success('Booking updated'),
    onError: async (e) => toast.error(await parseError(e)),
  });
}

export function useCancelBooking(bookingCode: string) {
  return useMutation({
    mutationFn: () => cancelBookingByCode(bookingCode),
    onSuccess: () => toast.success('Booking cancelled'),
    onError: async (e) => toast.error(await parseError(e)),
  });
}

export function useMyBookings(params?: { page?: number; size?: number }) {
  const page = params?.page ?? 1;
  const size = params?.size ?? 10;
  return useQuery({ 
    queryKey: ['my-bookings', page, size], 
    queryFn: () => getMyRegisteredBookings({ page, size }) 
  });
}

export function useRejectBooking(appCode?: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (bookingCode: string) => rejectBookingByCode(bookingCode),
    onSuccess: async () => {
      toast.success('Booking rejected');
      // Refresh related pages
      await Promise.all([
        // We cannot know which exact booking query keys exist; invalidate broad lists
        qc.invalidateQueries({ queryKey: ['appointment-users', appCode] }),
      ]);
    },
    onError: async (e) => toast.error(await parseError(e)),
  });
}

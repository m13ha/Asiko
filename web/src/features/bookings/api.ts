import type { BookingsApi, RequestsBookingRequest } from '@appointment-master/api-client';
import { bookingsApi } from '@/services/api';

export const bookingsClient: BookingsApi = bookingsApi;

export function getAvailableSlots(appCode: string, params?: { page?: number; size?: number }) {
  return bookingsClient.getAvailableSlots({ appCode, ...params });
}

export function getAvailableSlotsByDay(appCode: string, date: string, params?: { page?: number; size?: number }) {
  return bookingsClient.getAvailableSlotsByDay({ appCode, date, ...params });
}

export function bookGuest(booking: RequestsBookingRequest) {
  return bookingsClient.bookGuestAppointment({ booking });
}

export function bookRegistered(booking: RequestsBookingRequest) {
  return bookingsClient.bookRegisteredUserAppointment({ booking });
}

export function getBookingByCode(bookingCode: string) {
  return bookingsClient.getBookingByCode({ bookingCode });
}

export function updateBookingByCode(bookingCode: string, booking: RequestsBookingRequest) {
  return bookingsClient.updateBookingByCode({ bookingCode, booking });
}

export function cancelBookingByCode(bookingCode: string) {
  return bookingsClient.cancelBookingByCode({ bookingCode });
}

export function getMyRegisteredBookings(params?: { 
  page?: number; 
  size?: number;
  status?: string | string[];
}) {
  return bookingsClient.getUserRegisteredBookings(params);
}

export function rejectBookingByCode(bookingCode: string) {
  return bookingsClient.rejectBookingByCode({ bookingCode });
}

export function getAvailableDates(appCode: string) {
  return bookingsClient.getAvailableDates({ appCode });
}

export interface BookingRequest {
  appCode: string;
  name: string;
  email: string;
  phone: string;
  date: string;
  startTime: string;
  endTime: string;
  attendeeCount: number;
  description?: string;
}

export interface Booking {
  id: string;
  appointmentId?: string;
  userId?: string;
  name: string;
  email: string;
  phone: string;
  date: string;
  startTime: string;
  endTime: string;
  available: boolean;
  attendeeCount: number;
  createdAt: string;
  updatedAt: string;
  appCode: string;
  booking_code: string; // <-- Add this property for booking code
  description?: string;
}

// Default export to satisfy Expo Router
export default {};

export interface BookingRequest {
  appCode: string;
  name: string;
  email: string;
  phone: string;
  date: string;
  startTime: string;
  endTime: string;
  attendeeCount: number;
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
}
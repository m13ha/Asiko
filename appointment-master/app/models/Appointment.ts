export interface AppointmentRequest {
  title: string;
  startTime: string;
  endTime: string;
  startDate: string;
  endDate: string;
  bookingDuration: number;
  type: string;
  maxAttendees: number;
}

export interface AppointmentResponse {
  id: string;
  title: string;
  startTime: string;
  endTime: string;
  startDate: string;
  endDate: string;
  bookingDuration: number;
  type: string;
  maxAttendees: number;
  appCode: string;
  createdAt: string;
  updatedAt: string;
}
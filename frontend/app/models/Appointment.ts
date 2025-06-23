export interface AppointmentRequest {
  title: string;
  start_time: string;
  end_time: string;
  start_date: string;
  end_date: string;
  booking_duration: number;
  type: "single" | "group";
  max_attendees: number;
  description?: string;
}

export interface AppointmentResponse {
  id: string;
  title: string;
  start_time: string;
  end_time: string;
  start_date: string;
  end_date: string;
  booking_duration: number;
  type: string;
  max_attendees: number;
  app_code: string;
  created_at: string;
  updated_at: string;
  description?: string;
}

// Default export to satisfy Expo Router
export default {};

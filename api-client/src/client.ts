/**
 * Appointment Master API Client
 * 
 * Convenience wrapper for the generated API client with common configurations
 * and utility methods.
 */

import {
  Configuration,
  AuthenticationApi,
  AppointmentsApi,
  BookingsApi,
  AnalyticsApi,
  RequestsLoginRequest,
  RequestsUserRequest,
  RequestsAppointmentRequest,
  RequestsBookingRequest
} from './index';

export interface ClientConfig {
  baseUrl?: string;
  token?: string;
  timeout?: number;
}

export class AppointmentMasterClient {
  private config: Configuration;
  public auth: AuthenticationApi;
  public appointments: AppointmentsApi;
  public bookings: BookingsApi;
  public analytics: AnalyticsApi;

  constructor(clientConfig: ClientConfig = {}) {
    this.config = new Configuration({
      basePath: clientConfig.baseUrl || 'http://localhost:8888',
      accessToken: clientConfig.token,
    });

    this.auth = new AuthenticationApi(this.config);
    this.appointments = new AppointmentsApi(this.config);
    this.bookings = new BookingsApi(this.config);
    this.analytics = new AnalyticsApi(this.config);
  }

  /**
   * Update the authentication token
   */
  setToken(token: string) {
    this.config.accessToken = token;
  }

  /**
   * Remove the authentication token
   */
  clearToken() {
    this.config.accessToken = undefined;
  }

  /**
   * Login and automatically set the token
   */
  async login(email: string, password: string) {
    const response = await this.auth.loginUser({
      requestsLoginRequest: { email, password }
    });
    
    if (response.token) {
      this.setToken(response.token);
    }
    
    return response;
  }

  /**
   * Logout and clear the token
   */
  async logout() {
    try {
      await this.auth.logoutUser();
    } finally {
      this.clearToken();
    }
  }

  /**
   * Get user analytics for a date range
   */
  async getAnalytics(startDate: string, endDate: string) {
    return await this.analytics.getUserAnalytics({ startDate, endDate });
  }

  /**
   * Create a new user account
   */
  async register(userData: RequestsUserRequest) {
    return await this.auth.createUser({
      requestsUserRequest: userData
    });
  }

  /**
   * Create a new appointment
   */
  async createAppointment(appointmentData: RequestsAppointmentRequest) {
    return await this.appointments.createAppointment({
      requestsAppointmentRequest: appointmentData
    });
  }

  /**
   * Book an appointment as a guest
   */
  async bookAsGuest(bookingData: RequestsBookingRequest) {
    return await this.bookings.bookGuestAppointment({
      requestsBookingRequest: bookingData
    });
  }

  /**
   * Book an appointment as a registered user
   */
  async bookAsUser(bookingData: RequestsBookingRequest) {
    return await this.bookings.bookRegisteredUserAppointment({
      requestsBookingRequest: bookingData
    });
  }

  /**
   * Get available slots for an appointment
   */
  async getAvailableSlots(appointmentCode: string) {
    return await this.bookings.getAvailableSlots({ id: appointmentCode });
  }

  /**
   * Get booking details by booking code
   */
  async getBooking(bookingCode: string) {
    return await this.bookings.getBookingByCode({ bookingCode });
  }

  /**
   * Cancel a booking
   */
  async cancelBooking(bookingCode: string) {
    return await this.bookings.cancelBookingByCode({ bookingCode });
  }

  /**
   * Get user's appointments
   */
  async getMyAppointments() {
    return await this.appointments.getMyAppointments();
  }

  /**
   * Get user's bookings
   */
  async getMyBookings() {
    return await this.bookings.getUserRegisteredBookings();
  }
}

// Export everything from the generated client
export * from './index';

// Export the convenience client as default
export default AppointmentMasterClient;
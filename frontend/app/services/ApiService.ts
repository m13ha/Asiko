import AsyncStorage from "@react-native-async-storage/async-storage";
import { AppointmentRequest, AppointmentResponse } from "../models/Appointment";
import { BookingRequest, Booking } from "../models/Booking";
import { User } from "../models/User";
import { EventEmitter } from "events";
import { jwtDecode } from "jwt-decode";

import API_CONFIG from "./ApiConfig";
const API_URL = API_CONFIG.API_URL;

interface JwtPayload {
  exp: number;
  [key: string]: any;
}

class ApiService extends EventEmitter {
  private token: string | null = null;
  private user: User | null = null;

  constructor() {
    super();
    this.loadToken();
  }

  private isTokenExpired(token: string): boolean {
    try {
      const decoded = jwtDecode<JwtPayload>(token);
      if (!decoded.exp) return true;
      const now = Math.floor(Date.now() / 1000);
      return decoded.exp < now;
    } catch {
      return true;
    }
  }

  private async loadToken() {
    try {
      this.token = await AsyncStorage.getItem("token");
      if (this.token && this.isTokenExpired(this.token)) {
        await this.clearToken();
        this.user = null;
        this.emit("authChange", false);
      } else {
        this.emit("authChange", this.isAuthenticated);
      }
    } catch (error) {
      console.error("Failed to load token", error);
    }
  }

  private async saveToken(token: string) {
    try {
      await AsyncStorage.setItem("token", token);
      this.token = token;
      this.emit("authChange", this.isAuthenticated);
    } catch (error) {
      console.error("Failed to save token", error);
    }
  }

  private async clearToken() {
    try {
      await AsyncStorage.removeItem("token");
      this.token = null;
      this.emit("authChange", this.isAuthenticated);
    } catch (error) {
      console.error("Failed to clear token", error);
    }
  }

  private async request<T>(
    endpoint: string,
    method: string = "GET",
    data?: any
  ): Promise<T> {
    const url = `${API_URL}${endpoint}`;
    const headers: HeadersInit = {
      "Content-Type": "application/json",
    };

    if (this.token) {
      if (this.isTokenExpired(this.token)) {
        await this.clearToken();
        this.user = null;
        this.emit("authChange", false);
        throw new Error("Session expired. Please log in again.");
      }
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    const config: RequestInit = {
      method,
      headers,
      body: data ? JSON.stringify(data) : undefined,
    };

    try {
      const response = await fetch(url, config);
      let responseData;

      // Try to parse JSON response, handle text responses as well
      const contentType = response.headers.get("content-type");
      if (contentType && contentType.includes("application/json")) {
        responseData = await response.json();
      } else {
        const text = await response.text();
        try {
          responseData = JSON.parse(text);
        } catch (e) {
          responseData = { message: text };
        }
      }

      if (!response.ok) {
        console.error(`API Error (${response.status}):`, responseData);

        // If unauthorized or forbidden, clear token and force logout
        if (
          (response.status === 401 || response.status === 403) &&
          this.token
        ) {
          await this.clearToken();
          this.user = null;
          this.emit("authChange", false);
          throw new Error(
            "Session expired or unauthorized. Please log in again."
          );
        }

        // Format error message based on response structure
        let errorMessage = "API request failed";
        if (responseData.errors && Array.isArray(responseData.errors)) {
          errorMessage = responseData.errors
            .map((err: any) => `${err.field}: ${err.message}`)
            .join(", ");
        } else if (responseData.message) {
          errorMessage = responseData.message;
        } else if (typeof responseData === "string") {
          errorMessage = responseData;
        }
        throw new Error(errorMessage);
      }

      return responseData;
    } catch (error) {
      console.error(`API Request failed for ${method} ${endpoint}:`, error);
      throw error;
    }
  }

  // Custom method for auth change subscription
  subscribe(listener: (isAuthenticated: boolean) => void): () => void {
    this.addListener("authChange", listener);
    return () => {
      this.removeListener("authChange", listener);
    };
  }

  get isAuthenticated(): boolean {
    return !!this.token;
  }

  getUser(): User | null {
    return this.user;
  }

  async login(email: string, password: string): Promise<void> {
    const response = await this.request<{ token: string; user: User }>(
      "/login",
      "POST",
      {
        email,
        password,
      }
    );

    await this.saveToken(response.token);
    this.user = response.user;
  }

  async logout(): Promise<void> {
    await this.request("/logout", "POST");
    await this.clearToken();
    this.user = null;
  }

  async signup(
    name: string,
    email: string,
    phoneNumber: string,
    password: string
  ): Promise<void> {
    try {
      await this.request("/users", "POST", {
        name,
        email,
        phone_number: phoneNumber, // Use the correct field name expected by backend
        password,
      });
    } catch (error) {
      console.error("Signup error details:", error);
      // Rethrow the error for handling in the UI
      throw error;
    }
  }

  async createAppointment(
    data: AppointmentRequest
  ): Promise<AppointmentResponse> {
    return this.request<AppointmentResponse>("/appointments", "POST", data);
  }

  async getAppointments(): Promise<AppointmentResponse[]> {
    return this.request<AppointmentResponse[]>("/appointments/my");
  }

  // Add missing method that's being called in AppointmentsScreen
  async getUserAppointments(): Promise<AppointmentResponse[]> {
    return this.getAppointments();
  }

  async getAppointmentById(
    appointmentId: string
  ): Promise<AppointmentResponse> {
    // For now, we'll get all appointments and filter by ID
    // This should be replaced with a proper endpoint when available
    const appointments = await this.getAppointments();
    const appointment = appointments.find((app) => app.id === appointmentId);
    if (!appointment) {
      throw new Error("Appointment not found");
    }
    return appointment;
  }

  async bookGuestAppointment(data: BookingRequest): Promise<Booking> {
    return this.request<Booking>("/appointments/book", "POST", data);
  }

  async bookRegisteredUserAppointment(data: BookingRequest): Promise<Booking> {
    return this.request<Booking>("/appointments/book/registered", "POST", data);
  }

  async getAvailableSlots(appointmentId: string): Promise<Booking[]> {
    return this.request<Booking[]>(`/appointments/slots/${appointmentId}`);
  }

  async getUserBookings(): Promise<Booking[]> {
    return this.request<Booking[]>("/appointments/registered");
  }

  async getBookingByCode(bookingCode: string): Promise<Booking> {
    return this.request<Booking>(`/bookings/${bookingCode}`);
  }

  async updateBooking(
    bookingCode: string,
    data: Partial<BookingRequest>
  ): Promise<Booking> {
    return this.request<Booking>(`/bookings/${bookingCode}`, "PUT", data);
  }

  async cancelBooking(bookingCode: string): Promise<void> {
    await this.request(`/bookings/${bookingCode}`, "DELETE");
  }
}

export default new ApiService();

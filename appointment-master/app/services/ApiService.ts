import AsyncStorage from '@react-native-async-storage/async-storage';
import { User, LoginRequest, UserRequest } from '../models/User';
import { AppointmentRequest, AppointmentResponse } from '../models/Appointment';
import { Booking, BookingRequest } from '../models/Booking';

class ApiService {
  private token: string | null = null;
  private user: User | null = null;
  private baseUrl = 'http://127.0.0.1:8080';
  private listeners: (() => void)[] = [];

  constructor() {
    this.loadTokenAndUser();
  }

  get isAuthenticated(): boolean {
    return this.token !== null;
  }

  addListener(listener: () => void) {
    this.listeners.push(listener);
    return () => {
      this.listeners = this.listeners.filter(l => l !== listener);
    };
  }

  private notifyListeners() {
    this.listeners.forEach(listener => listener());
  }

  private async loadTokenAndUser() {
    try {
      const token = await AsyncStorage.getItem('token');
      const userJson = await AsyncStorage.getItem('user');
      
      if (token && userJson) {
        this.token = token;
        this.user = JSON.parse(userJson);
        this.notifyListeners();
      } else if (token) {
        this.token = token;
        await this.fetchUserDetails();
      }
    } catch (error) {
      console.error('Error loading token and user:', error);
    }
  }

  private async saveTokenAndUser(token: string, user: User) {
    try {
      await AsyncStorage.setItem('token', token);
      await AsyncStorage.setItem('user', JSON.stringify(user));
      this.token = token;
      this.user = user;
      this.notifyListeners();
    } catch (error) {
      console.error('Error saving token and user:', error);
    }
  }

  private async fetchUserDetails() {
    try {
      const response = await fetch(`${this.baseUrl}/me`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.token}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        this.user = data;
        await AsyncStorage.setItem('user', JSON.stringify(this.user));
        this.notifyListeners();
      } else {
        throw new Error(`Failed to fetch user details: ${await response.text()}`);
      }
    } catch (error) {
      console.error('Error fetching user details:', error);
      throw error;
    }
  }

  async login(email: string, password: string) {
    try {
      const response = await fetch(`${this.baseUrl}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        const data = await response.json();
        const token = data.token;
        
        if (data.user) {
          await this.saveTokenAndUser(token, data.user);
        } else {
          await this.saveTokenAndUser(token, {
            id: '',
            name: '',
            email,
            phone: '',
          });
          await this.fetchUserDetails();
        }
        return true;
      } else {
        throw new Error(`Login failed: ${await response.text()}`);
      }
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  async signup(name: string, email: string, phone: string, password: string) {
    try {
      const response = await fetch(`${this.baseUrl}/users`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name,
          email,
          phone_number: phone,
          password,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        await this.saveTokenAndUser(data.token, data);
        return true;
      } else {
        throw new Error(`Signup failed: ${await response.text()}`);
      }
    } catch (error) {
      console.error('Signup error:', error);
      throw error;
    }
  }

  async logout() {
    try {
      await AsyncStorage.removeItem('token');
      await AsyncStorage.removeItem('user');
      this.token = null;
      this.user = null;
      this.notifyListeners();
    } catch (error) {
      console.error('Logout error:', error);
      throw error;
    }
  }

  async getUserAppointments(): Promise<AppointmentResponse[]> {
    if (!this.token) throw new Error('Not authenticated');
    
    try {
      const response = await fetch(`${this.baseUrl}/appointments/my`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.token}`,
        },
      });

      if (response.ok) {
        return await response.json();
      } else {
        throw new Error(`Failed to fetch appointments: ${await response.text()}`);
      }
    } catch (error) {
      console.error('Error fetching appointments:', error);
      throw error;
    }
  }

  async getAppointmentBookings(appointmentId: string): Promise<Booking[]> {
    if (!this.token) throw new Error('Not authenticated');
    
    try {
      const response = await fetch(`${this.baseUrl}/appointments/users/${appointmentId}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.token}`,
        },
      });

      if (response.ok) {
        return await response.json();
      } else {
        throw new Error(`Failed to fetch bookings: ${await response.text()}`);
      }
    } catch (error) {
      console.error('Error fetching bookings:', error);
      throw error;
    }
  }

  async getAvailableSlots(appointmentId: string): Promise<Booking[]> {
    try {
      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
      };
      
      if (this.token) {
        headers['Authorization'] = `Bearer ${this.token}`;
      }
      
      const response = await fetch(`${this.baseUrl}/appointments/slots/${appointmentId}`, {
        headers,
      });

      if (response.ok) {
        return await response.json();
      } else {
        throw new Error(`Failed to fetch available slots: ${await response.text()}`);
      }
    } catch (error) {
      console.error('Error fetching available slots:', error);
      throw error;
    }
  }

  async bookGuestAppointment(booking: BookingRequest): Promise<Booking> {
    try {
      const response = await fetch(`${this.baseUrl}/appointments/book`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(booking),
      });

      if (response.ok) {
        return await response.json();
      } else {
        throw new Error(`Failed to book guest appointment: ${await response.text()}`);
      }
    } catch (error) {
      console.error('Error booking guest appointment:', error);
      throw error;
    }
  }

  async bookRegisteredUserAppointment(booking: BookingRequest): Promise<Booking> {
    if (!this.token) throw new Error('Not authenticated');
    
    try {
      const response = await fetch(`${this.baseUrl}/appointments/book/registered`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.token}`,
        },
        body: JSON.stringify(booking),
      });

      if (response.ok) {
        return await response.json();
      } else {
        throw new Error(`Failed to book registered user appointment: ${await response.text()}`);
      }
    } catch (error) {
      console.error('Error booking registered user appointment:', error);
      throw error;
    }
  }

  async createAppointment(appointment: AppointmentRequest): Promise<AppointmentResponse> {
    if (!this.token) throw new Error('Not authenticated');
    
    try {
      const response = await fetch(`${this.baseUrl}/appointments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.token}`,
        },
        body: JSON.stringify(appointment),
      });

      if (response.ok) {
        return await response.json();
      } else {
        throw new Error(`Failed to create appointment: ${await response.text()}`);
      }
    } catch (error) {
      console.error('Error creating appointment:', error);
      throw error;
    }
  }

  getUser(): User | null {
    return this.user;
  }
}

// Create a singleton instance
const apiService = new ApiService();
export default apiService;
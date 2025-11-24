# Asiko API Client

A TypeScript/JavaScript client library for the Asiko API, automatically generated from OpenAPI specifications.

## Installation

```bash
npm install @appointment-master/api-client
# or
yarn add @appointment-master/api-client
```

## Quick Start

```typescript
import { 
  AuthenticationApi, 
  AppointmentsApi, 
  BookingsApi,
  Configuration 
} from '@appointment-master/api-client';

// Configure the API client
const config = new Configuration({
  basePath: 'http://localhost:8888',
  // Inject Bearer token for authenticated calls
  apiKey: () => `Bearer ${localStorage.getItem('authToken') ?? ''}`,
});

// Initialize API clients
const authApi = new AuthenticationApi(config);
const appointmentsApi = new AppointmentsApi(config);
const bookingsApi = new BookingsApi(config);
```

## Authentication

### Login

```typescript
import { AuthenticationApi, RequestsLoginRequest } from '@appointment-master/api-client';

const authApi = new AuthenticationApi(config);

const loginRequest: RequestsLoginRequest = {
  email: 'user@example.com',
  password: 'password123'
};

try {
  const response = await authApi.loginUser({ requestsLoginRequest: loginRequest });
  const { token, refreshToken, expiresIn, user } = response;
  
  // Store token for future requests
  localStorage.setItem('authToken', token);
  localStorage.setItem('refreshToken', refreshToken);
  
  // Update configuration with token
  config.apiKey = () => `Bearer ${token}`;
} catch (error) {
  console.error('Login failed:', error);
}
```

### Refresh Access Token

```typescript
const refreshToken = localStorage.getItem('refreshToken');
if (refreshToken) {
  const refreshed = await authApi.refreshToken({
    refresh: { refreshToken }
  });
  localStorage.setItem('authToken', refreshed.token);
  localStorage.setItem('refreshToken', refreshed.refreshToken);
  config.apiKey = () => `Bearer ${refreshed.token}`;
}
```

### Logout

```typescript
try {
  await authApi.logoutUser();
  localStorage.removeItem('authToken');
} catch (error) {
  console.error('Logout failed:', error);
}
```

## User Management

### Create User

```typescript
import { RequestsUserRequest } from '@appointment-master/api-client';

const userRequest: RequestsUserRequest = {
  name: 'John Doe',
  email: 'john@example.com',
  password: 'securepassword123',
  phoneNumber: '+1234567890'
};

try {
  const user = await authApi.createUser({ requestsUserRequest: userRequest });
  console.log('User created:', user);
} catch (error) {
  console.error('User creation failed:', error);
}
```

## Appointments

### Create Appointment

```typescript
import { RequestsAppointmentRequest, EntitiesAppointmentType } from '@appointment-master/api-client';

const appointmentRequest: RequestsAppointmentRequest = {
  title: 'Consultation Session',
  startTime: new Date('2025-01-15T09:00:00Z'),
  endTime: new Date('2025-01-15T17:00:00Z'),
  startDate: new Date('2025-01-15T00:00:00Z'),
  endDate: new Date('2025-01-20T00:00:00Z'),
  bookingDuration: 60, // minutes
  type: EntitiesAppointmentType.Single,
  maxAttendees: 1,
  description: 'Professional consultation'
};

try {
  const appointment = await appointmentsApi.createAppointment({ 
    requestsAppointmentRequest: appointmentRequest 
  });
  console.log('Appointment created:', appointment);
} catch (error) {
  console.error('Appointment creation failed:', error);
}
```

### Get My Appointments

```typescript
try {
  const myAppointments = await appointmentsApi.getMyAppointments();
  console.log('My appointments:', myAppointments);
} catch (error) {
  console.error('Failed to fetch appointments:', error);
}
```

## Bookings

### Book Appointment (Guest)

```typescript
import { RequestsBookingRequest } from '@appointment-master/api-client';

const bookingRequest: RequestsBookingRequest = {
  appCode: 'AP-ABC123',
  startTime: new Date('2025-01-15T10:00:00Z'),
  endTime: new Date('2025-01-15T11:00:00Z'),
  date: new Date('2025-01-15T00:00:00Z'),
  name: 'Jane Smith',
  email: 'jane@example.com',
  phone: '+1987654321',
  attendeeCount: 1,
  description: 'Initial consultation'
};

try {
  const booking = await bookingsApi.bookGuestAppointment({ 
    requestsBookingRequest: bookingRequest 
  });
  console.log('Booking created:', booking);
} catch (error) {
  console.error('Booking failed:', error);
}
```

### Book Appointment (Registered User)

```typescript
// Requires authentication token in configuration
const bookingRequest: RequestsBookingRequest = {
  appCode: 'AP-ABC123',
  startTime: new Date('2025-01-15T10:00:00Z'),
  endTime: new Date('2025-01-15T11:00:00Z'),
  date: new Date('2025-01-15T00:00:00Z'),
  attendeeCount: 1,
  description: 'Follow-up consultation'
};

try {
  const booking = await bookingsApi.bookRegisteredUserAppointment({ 
    requestsBookingRequest: bookingRequest 
  });
  console.log('Booking created:', booking);
} catch (error) {
  console.error('Booking failed:', error);
}
```

### Get Available Slots

```typescript
try {
  const slots = await bookingsApi.getAvailableSlots({ id: 'AP-ABC123' });
  console.log('Available slots:', slots);
} catch (error) {
  console.error('Failed to fetch slots:', error);
}
```

### Get Booking by Code

```typescript
try {
  const booking = await bookingsApi.getBookingByCode({ bookingCode: 'BK-XYZ789' });
  console.log('Booking details:', booking);
} catch (error) {
  console.error('Booking not found:', error);
}
```

### Cancel Booking

```typescript
try {
  const cancelledBooking = await bookingsApi.cancelBookingByCode({ 
    bookingCode: 'BK-XYZ789' 
  });
  console.log('Booking cancelled:', cancelledBooking);
} catch (error) {
  console.error('Cancellation failed:', error);
}
```

## Error Handling

The client throws errors that can be caught and handled:

```typescript
import { ResponseError } from '@appointment-master/api-client';

try {
  await authApi.loginUser({ requestsLoginRequest: loginData });
} catch (error) {
  if (error instanceof ResponseError) {
    console.error('API Error:', error.response.status, error.response.statusText);
    
    // Get error details
    const errorBody = await error.response.json();
    console.error('Error details:', errorBody);
  } else {
    console.error('Network or other error:', error);
  }
}
```

## Configuration Options

```typescript
const config = new Configuration({
  basePath: 'http://localhost:8888', // API base URL
  accessToken: 'jwt-token', // JWT token for authentication
  fetchApi: fetch, // Custom fetch implementation
  middleware: [
    // Custom middleware for requests/responses
    {
      pre: async (context) => {
        console.log('Request:', context.url);
        return context;
      },
      post: async (context) => {
        console.log('Response:', context.response.status);
        return context.response;
      }
    }
  ]
});
```

## React Hook Example

```typescript
import { useState, useEffect } from 'react';
import { AuthenticationApi, Configuration } from '@appointment-master/api-client';

export function useAuth() {
  const [token, setToken] = useState<string | null>(
    localStorage.getItem('authToken')
  );
  const [user, setUser] = useState(null);

  const config = new Configuration({
    basePath: process.env.REACT_APP_API_URL || 'http://localhost:8888',
    accessToken: token || undefined
  });

  const authApi = new AuthenticationApi(config);

  const login = async (email: string, password: string) => {
    try {
      const response = await authApi.loginUser({
        requestsLoginRequest: { email, password }
      });
      
      setToken(response.token);
      setUser(response.user);
      localStorage.setItem('authToken', response.token);
      
      return response;
    } catch (error) {
      throw error;
    }
  };

  const logout = async () => {
    try {
      await authApi.logoutUser();
    } finally {
      setToken(null);
      setUser(null);
      localStorage.removeItem('authToken');
    }
  };

  return { token, user, login, logout, config };
}
```

## Development

### Building the Client

```bash
# Install dependencies
npm install

# Build TypeScript to JavaScript
npm run build

# Watch for changes during development
npm run dev
```

### Regenerating the Client

When the API changes, regenerate the client:

```bash
# From the backend directory
cd backend
npm run client-gen
```

## API Reference

### Available APIs

- **AuthenticationApi**: User login, logout, registration
- **AppointmentsApi**: Create and manage appointments
- **BookingsApi**: Book appointments, manage bookings, get availability

### Models

All TypeScript interfaces are available for import:

- `RequestsLoginRequest`, `RequestsUserRequest`, `RequestsAppointmentRequest`, `RequestsBookingRequest`
- `EntitiesAppointment`, `EntitiesBooking`, `EntitiesAppointmentType`
- `ResponsesUserResponse`, `ResponsesAppointmentResponse`, `ResponsesPaginatedResponse`

## Support

For issues and questions:
- Check the [API documentation](http://localhost:8888/swagger/index.html)
- Review the generated TypeScript interfaces
- Check network requests in browser dev tools

## License

MIT License - see LICENSE file for details.

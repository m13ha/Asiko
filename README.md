# Appointment Booking App (In Development)

## Overview

The **Appointment Booking App** is designed to provide an easy-to-use appointment creation and booking management tool tailored for small businesses and companies. Built with simplicity and efficiency in mind, this application allows users to create appointments, manage bookings, and query available slots seamlessly. It supports both registered users and guests, making it versatile for various use cases.

This project is written in **Golang** using the **GORM** ORM and **Chi router**, providing a robust backend API for appointment management. The current implementation serves as a solid foundation, with room for future enhancements like payment integration, notifications, and booking approval systems.

## Features

### Current Features

- **User Registration and Authentication**
  - Users can create accounts with usernames, emails, and passwords.
  - Secure login/logout functionality using JWT-based authentication.
- **Appointment Creation**
  - Registered users can create appointments with customizable start/end times, dates, and booking durations.
  - Supports both single and group appointment types with configurable maximum attendees.
  - Automatically generates time slots based on booking duration.
- **Booking Management**
  - Registered users and guests can book appointments.
  - Guests must provide a name and contact information (email or phone).
  - Registered users’ details are auto-filled from their profiles.
  - Users can retrieve appointments they’ve created or booked.
- **Slot Availability**
  - API endpoint to query available slots for any appointment.
- **Data Persistence**
  - Uses PostgreSQL with GORM for reliable data storage.
  - Soft deletes implemented for data recovery.

### Planned Enhancements

- **Payment Integration**: Add support for processing payments for bookings.
- **Notifications**: Implement email/SMS notifications for booking confirmations, reminders, and updates.
- **Booking Approval System**: Allow business owners to approve or reject bookings manually.
- **Calendar Integration**: Sync appointments with external calendars (e.g., Google Calendar).
- **Advanced Management Tools**: Add features like booking cancellation, rescheduling, and recurring appointments.

## Project Structure

```
root/
├── API/           # API handlers for routing and request handling
├── Models/        # Data models and structs
├── Services/      # Business logic and database interactions
├── Database/      # Database connection and migration logic
├── Utils/         # Utility functions (e.g., validation, code generation)
├── main.go        # Entry point and server setup
└── auth.go        # Authentication logic (moved into API package)
```

## Tech Stack

- **Backend**: Golang
- **Framework**: Chi Router
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **Environment**: Managed with `.env` file using `godotenv`

## Getting Started

### Prerequisites

- Go 1.20 or later
- PostgreSQL
- Git

### Installation

1. **Clone the Repository**

   ```bash
   git clone <repository-url>
   cd appointment-booking-app
   ```

2. **Set Up Environment Variables**
   Create a `.env` file in the root directory with the following:

   ```env
   PORT=8080
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=appointmentdb
   DB_SSLMODE=disable
   JWT_SECRET_KEY=your_secret_key_here
   ```

3. **Install Dependencies**

   ```bash
   go mod tidy
   ```

4. **Run Database Migrations**
   The app automatically migrates the database schema on startup (drops and recreates tables in development mode).

5. **Start the Server**
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:8080` (or the port specified in `.env`).

### API Endpoints

- **Public Routes**

  - `POST /login` - User login
  - `POST /logout` - User logout (client-side token discard)
  - `POST /users` - Create a new user
  - `POST /appointments/book` - Book an appointment (guest or registered)
  - `GET /appointments/{id}/slots` - Get available slots for an appointment

- **Protected Routes** (requires JWT token in `Authorization: Bearer <token>` header)
  - `POST /appointments` - Create a new appointment
  - `GET /appointments/{id}/users` - Get all bookings for an appointment
  - `GET /appointments/my` - Get appointments created by the user
  - `GET /appointments/registered` - Get appointments booked by the user

## Usage Examples

### Create a User

```bash
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{"name": "John Doe", "email": "john@example.com", "password": "securepass123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{"email": "john@example.com", "password": "securepass123"}'
```

### Book an Appointment (Guest)

```bash
curl -X POST http://localhost:8080/appointments/book \
-H "Content-Type: application/json" \
-d '{"appointment_id": "<uuid>", "start_time": "2025-03-22T10:00:00Z", "end_time": "2025-03-22T10:30:00Z", "date": "2025-03-22T00:00:00Z", "name": "Jane Doe", "email": "jane@example.com", "attendee_count": 1}'
```

### Book an Appointment (Registered User)

```bash
curl -X POST http://localhost:8080/appointments/book \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <token>" \
-d '{"appointment_id": "<uuid>", "start_time": "2025-03-22T10:00:00Z", "end_time": "2025-03-22T10:30:00Z", "date": "2025-03-22T00:00:00Z", "attendee_count": 1}'
```

## Future Development

This project is designed for extensibility. Potential upgrades include:

- **Frontend**: Build a user-friendly interface using Flutter (see below).
- **Payment Gateway**: Integrate Stripe or PayPal for paid bookings.
- **Notifications**: Use Twilio for SMS or SendGrid for email notifications.
- **Approval Workflow**: Add a status field to bookings and endpoints for approval/rejection.

## Contributing

Contributions are welcome! Please fork the repository, create a feature branch, and submit a pull request with your changes.

## License

This project is licensed under the MIT License.

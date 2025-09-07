# Appointment Master Backend API  

## 🚀 Executive Summary  
Appointment Master is a **scalable booking and ticketing API** designed to make scheduling **fast, flexible, and user-friendly**. Built with **Go (Gin) + PostgreSQL**, it supports **both guest and registered users**, automatic slot generation, and **real-time availability checks**.  

### Highlights  
- 🔹 REST API for appointment booking & event ticketing  
- 🔹 Built with **Go, Gin, GORM, PostgreSQL**  
- 🔹 **JWT authentication** with secure password hashing  
- 🔹 Supports **single, group, and party-style bookings**  
- 🔹 **Automatic slot generation** with live availability  
- 🔹 **Dockerized deployment** + Swagger/OpenAPI docs  
- 🔹 Designed with **Clean Architecture** for scalability  

---

## 📖 Overview  
The **Appointment Master Backend** is a production-ready API that simplifies appointment management for **small businesses, service providers, and event organizers**. Whether it’s a barber managing daily clients, a clinic scheduling consultations, or a teacher setting up office hours, Appointment Master provides an **easy, reliable way to handle bookings and tickets**.  

Core benefits:  
- 💡 **Flexible Appointment Types** – one-on-one, group sessions, or large event bookings  
- 🔒 **Secure User Management** – JWT auth, bcrypt password hashing, and unique IDs  
- ⏱ **Smart Scheduling** – real-time availability with database locking to prevent double-bookings  
- 🛠 **Developer Friendly** – Docker setup, Swagger documentation, and structured logging  

---

## Architecture  

### Clean Architecture Implementation
```

┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Layer     │    │  Service Layer  │    │Repository Layer │
│   (Handlers)    │───▶│ (Business Logic)│───▶│ (Data Access)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
│                       │                       │
▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Middleware    │    │     Models      │    │    Database     │
│ (Auth, CORS,    │    │ (Entities,      │    │  (PostgreSQL)   │
│  Logging)       │    │  Requests)      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘

```

### Project Structure
```

backend/
├── api/                    # HTTP handlers and routing
│   ├── appointment.go      # Appointment endpoints
│   ├── auth.go             # Authentication endpoints
│   ├── booking.go          # Booking endpoints
│   ├── handlers.go         # Route registration
│   └── user.go             # User management endpoints
├── db/                     # Database configuration
│   └── db.go               # Connection, migration, pooling
├── middleware/             # HTTP middleware
│   ├── auth.go             # JWT authentication
│   ├── cors.go             # Cross-origin resource sharing
│   └── logger.go           # Request logging
├── models/                 # Data models
│   ├── entities/           # Database entities
│   ├── requests/           # API request models
│   └── responses/          # API response models
├── repository/             # Data access layer
│   ├── appointment\_repository.go
│   ├── booking\_repository.go
│   └── user\_repository.go
├── services/               # Business logic layer
│   ├── appointment.go      # Appointment business logic
│   ├── booking.go          # Booking business logic
│   └── user.go             # User business logic
├── utils/                  # Utility functions
│   ├── appcode.go          # Code generation
│   ├── email.go            # Email utilities
│   ├── string.go           # String utilities
│   ├── time.go             # Time utilities
│   └── validate.go         # Input validation
└── main.go                 # Application entry point

````

---

## Core Features  

### Authentication & User Management
- **User Registration**: Create accounts with name, email, password, and optional phone number  
- **JWT Authentication**: Secure token-based authentication with 24-hour expiration  
- **Email Normalization**: Automatic email case normalization for consistency  
- **Password Security**: bcrypt hashing with salt for secure password storage  
- **Duplicate Prevention**: Email and phone number uniqueness validation  

### Appointment Management
- **Flexible Appointment Types**:
  - **Single**: One-on-one appointments with automatic slot generation  
  - **Group**: Multiple attendees with configurable capacity limits  
  - **Party**: Event-style bookings with shared capacity management  
- **Smart Slot Generation**: Automatic time slot creation based on booking duration  
- **Date Range Support**: Multi-day appointment scheduling  
- **Unique App Codes**: Auto-generated appointment codes (AP-XXXXX format)  

### Booking System
- **Dual Booking Support**: Both registered users and guest bookings  
- **Real-time Availability**: Live slot availability checking with database locking  
- **Capacity Management**: Automatic attendee count tracking and validation  
- **Booking Codes**: Unique booking identifiers (BK-XXXXX format) for easy reference  
- **Booking Operations**: Create, update, cancel, and retrieve bookings  
- **Status Tracking**: Active/cancelled booking status management  

### Data Management
- **PostgreSQL Database**: Robust relational database with ACID compliance  
- **Soft Deletes**: Data recovery capability with GORM soft delete implementation  
- **UUID Primary Keys**: Distributed-system-friendly unique identifiers  
- **Database Migrations**: Automatic schema management and updates  
- **Connection Pooling**: Optimized database connection management  

---

## Tech Stack  

- **Language**: Go 1.24  
- **Framework**: Gin (HTTP web framework)  
- **ORM**: GORM v1.25.12  
- **Database**: PostgreSQL 17  
- **Authentication**: JWT with HS256 signing  
- **Validation**: go-playground/validator/v10  
- **Logging**: zerolog (structured logging)  
- **Testing**: testify (mocking and assertions)  
- **Containerization**: Docker & Docker Compose  
- **Documentation**: Swagger/OpenAPI  

---

## API Endpoints  

### Public Endpoints (No Authentication Required)  

#### Authentication
- `POST /login` – User authentication  
- `POST /logout` – User logout (client-side token invalidation)  
- `POST /users` – User registration  

#### Booking Management
- `POST /appointments/book` – Book appointment (guest or registered user)  
- `GET /appointments/slots/:id` – Get available slots for appointment  
- `GET /appointments/slots/:id/by-day` – Get available slots by specific day  
- `GET /bookings/:booking_code` – Get booking details by code  
- `PUT /bookings/:booking_code` – Update booking by code  
- `DELETE /bookings/:booking_code` – Cancel booking by code  

### Protected Endpoints (JWT Authentication Required)  

#### Appointment Management
- `POST /appointments` – Create new appointment  
- `GET /appointments/my` – Get user's created appointments  
- `GET /appointments/users/:id` – Get all bookings for appointment  

#### User Bookings
- `GET /appointments/registered` – Get user's booked appointments  
- `POST /appointments/book/registered` – Book appointment as registered user  

---

## Data Models  

### User Entity
```go
type User struct {
    ID             uuid.UUID      `json:"id"`
    Name           string         `json:"name"`
    Email          string         `json:"email"`
    PhoneNumber    string         `json:"phone_number"`
    HashedPassword string         `json:"-"`
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
    DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty"`
}
````

### Appointment Entity

```go
type Appointment struct {
    ID              uuid.UUID       `json:"id"`
    Title           string          `json:"title"`
    StartTime       time.Time       `json:"start_time"`
    EndTime         time.Time       `json:"end_time"`
    StartDate       time.Time       `json:"start_date"`
    EndDate         time.Time       `json:"end_date"`
    BookingDuration int             `json:"booking_duration"` // minutes
    MaxAttendees    int             `json:"max_attendees"`
    Type            AppointmentType `json:"type"` // single, group, party
    OwnerID         uuid.UUID       `json:"owner_id"`
    AppCode         string          `json:"app_code"`
    Description     string          `json:"description"`
    AttendeesBooked int             `json:"attendees_booked"`
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at"`
    DeletedAt       gorm.DeletedAt  `json:"deleted_at,omitempty"`
}
```

### Booking Entity

```go
type Booking struct {
    ID                  uuid.UUID      `json:"id"`
    AppointmentID       uuid.UUID      `json:"appointment_id"`
    AppCode             string         `json:"app_code"`
    UserID              *uuid.UUID     `json:"user_id"` // null for guest bookings
    Name                string         `json:"name"`
    Email               string         `json:"email"`
    Phone               string         `json:"phone"`
    Date                time.Time      `json:"date"`
    StartTime           time.Time      `json:"start_time"`
    EndTime             time.Time      `json:"end_time"`
    Available           bool           `json:"available"`
    AttendeeCount       int            `json:"attendee_count"`
    BookingCode         string         `json:"booking_code"`
    Status              string         `json:"status"` // active, cancelled
    Description         string         `json:"description"`
    CreatedAt           time.Time      `json:"created_at"`
    UpdatedAt           time.Time      `json:"updated_at"`
    DeletedAt           gorm.DeletedAt `json:"deleted_at,omitempty"`
}
```

---

## How It Works

### 1. User Registration & Authentication Flow

```
1. User submits registration data (name, email, password, phone)
2. System validates input and checks for duplicates
3. Password is hashed using bcrypt
4. User record is created in database
5. For login: credentials are verified and JWT token is issued
```

### 2. Appointment Creation Flow

```
1. Authenticated user submits appointment details
2. System validates appointment data and time ranges
3. Appointment is created with unique app code (AP-XXXXX)
4. For single/group types: time slots are auto-generated
5. For party type: no slots generated (shared capacity model)
```

### 3. Booking Flow

#### Guest Booking

```
1. Guest provides appointment code and personal details
2. System validates appointment exists and has availability
3. For slot-based: finds and reserves specific time slot
4. For party-based: checks capacity and creates booking
5. Booking created with unique booking code (BK-XXXXX)
```

#### Registered User Booking

```
1. Authenticated user provides appointment code and preferences
2. System auto-fills user details from profile
3. Same availability checking and reservation process
4. Booking linked to user account for easy retrieval
```

### 4. Slot Management System

#### Single/Group Appointments

* Slots are pre-generated based on booking duration
* Each slot has availability flag and attendee count
* Booking marks slot as unavailable or increments count
* Cancellation restores availability

#### Party Appointments

* No pre-generated slots
* Uses shared capacity model with attendee tracking
* Real-time capacity checking with database locking
* Supports concurrent bookings up to max capacity

---

## Environment Configuration

### Required Environment Variables

```env
# Server Configuration
PORT=8888
ENV=development

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=password
DB_NAME=appointmentdb

# Authentication
JWT_SECRET_KEY=your_secret_key_here

# PostgreSQL Docker Configuration
POSTGRES_PASSWORD=password
POSTGRES_DB=postgres
```

---

## Development Setup

### Using Docker (Recommended)

```bash
# Clone repository
git clone <repository-url>
cd appointment-master/backend

# Start services
docker-compose up --build

# API available at http://localhost:8888
# Swagger docs at http://localhost:8888/swagger/index.html
```

### Local Development

```bash
# Install dependencies
go mod tidy

# Set up PostgreSQL database
createdb appointmentdb

# Run with hot reload
make dev

# Run tests
make test

# Build for production
make build
```

---

## Testing

### Test Coverage

* **Unit Tests**: Service layer business logic
* **Integration Tests**: API endpoint testing
* **Mock Testing**: Repository layer mocking
* **Validation Tests**: Input validation scenarios

### Running Tests

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific package tests
go test ./services -v
```

---

## API Usage Examples

### 1. User Registration

```bash
curl -X POST http://localhost:8888/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepass123",
    "phone_number": "+1234567890"
  }'
```

### 2. User Login

```bash
curl -X POST http://localhost:8888/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

### 3. Create Appointment

```bash
curl -X POST http://localhost:8888/appointments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "title": "Consultation Session",
    "start_time": "2025-01-15T09:00:00Z",
    "end_time": "2025-01-15T17:00:00Z",
    "start_date": "2025-01-15T00:00:00Z",
    "end_date": "2025-01-20T00:00:00Z",
    "booking_duration": 60,
    "type": "single",
    "max_attendees": 1,
    "description": "Professional consultation"
  }'
```

### 4. Guest Booking

```bash
curl -X POST http://localhost:8888/appointments/book \
  -H "Content-Type: application/json" \
  -d '{
    "app_code": "AP-ABC123",
    "start_time": "2025-01-15T10:00:00Z",
    "end_time": "2025-01-15T11:00:00Z",
    "date": "2025-01-15T00:00:00Z",
    "name": "Jane Smith",
    "email": "jane@example.com",
    "phone": "+1987654321",
    "attendee_count": 1,
    "description": "Initial consultation"
  }'
```

### 5. Get Available Slots

```bash
curl -X GET http://localhost:8888/appointments/slots/AP-ABC123 \
  -H "Content-Type: application/json"
```

### 6. Cancel Booking

```bash
curl -X DELETE http://localhost:8888/bookings/BK-XYZ789 \
  -H "Content-Type: application/json"
```

---

## Production Deployment

### Docker Production Setup

```yaml
# docker-compose.prod.yml
version: "3.8"
services:
  api:
    build: .
    environment:
      - ENV=production
      - JWT_SECRET_KEY=${JWT_SECRET_KEY}
    ports:
      - "8888:8080"
    depends_on:
      - postgres
    restart: always

  postgres:
    image: postgres:17-alpine
    environment:
      - POSTGRES_DB=${DB_NAME}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: always
```

### Security Considerations

* JWT tokens expire after 24 hours
* Passwords are hashed with bcrypt (cost 12)
* SQL injection prevention with parameterized queries
* Input validation on all endpoints
* CORS configuration for cross-origin requests
* Rate limiting recommended for production

---

## Monitoring & Observability

### Structured Logging

* Request/response logging with zerolog
* Error tracking with stack traces
* Performance metrics logging
* Database query logging in development

### Health Checks

* Database connectivity checks
* Application health endpoints
* Docker health check configuration

---

## 🔮 Future Roadmap

Appointment Master is designed with **scalability and extensibility** in mind. The following features can be added to evolve this project toward a SaaS-grade product:

### User Experience Enhancements

* 📧 **Email Notifications** – confirmations, cancellations, and reminders
* 📱 **SMS Integration** – instant text updates for bookings
* 📆 **Calendar Sync** – Google Calendar and Outlook integration

### Business Features

* 💳 **Payment Processing** – Stripe/PayPal integration for paid bookings
* 🔁 **Recurring Appointments** – repeat booking patterns for ongoing services
* ✅ **Approval Workflow** – manual booking approval for businesses that need control
* 📊 **Analytics Dashboard** – insights into booking trends and utilization

### Technical Improvements

* ⚡ **Caching Layer** – Redis for session management and faster queries
* ⛔ **Rate Limiting** – request throttling to prevent abuse
* 📈 **Metrics & Monitoring** – Prometheus and Grafana integration
* 🔎 **Distributed Tracing** – full request lifecycle tracking across services
* 🗄 **Database Optimization** – indexing and query tuning for scale

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines


* Follow Clean Architecture principles
* Write tests for new features
* Maintain code documentation
* Use meaningful commit messages
* Ensure all tests pass before PR submission

---

## License

MIT License – feel free to use, modify, and distribute with attribution.

```
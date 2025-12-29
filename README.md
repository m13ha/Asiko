# Asiko Platform

**Portfolio Project Notice**  
Asiko is a portfolio project and is free for anybody and everybody to use to manage appointments. The codebase is designed to be easy to use and easy to manage for developers and operators.

Asiko is a full-stack appointment and booking platform that helps solo operators and small teams publish availability, accept reservations, and keep owners and attendees in sync. The system combines a Go API, a React web client, and shared tooling in a deployable monorepo.

---

## The Problem It Solves

Service businesses often juggle spreadsheets, messaging threads, and multiple tools to answer four questions: *What slots are open? Who booked them? Did we confirm? Can we trust the data?* Asiko addresses those gaps by:

- Keeping a single source of truth for appointments, slots, bookings, and notifications.
- Supporting both owner workflows and guest booking flows within one API.
- Exposing consistent analytics and status transitions across the lifecycle.
- Providing a web UI that shares a generated API client to stay aligned with backend contracts.

---

## Architecture Overview

```
root
├─ backend/     # Go + Gin API
├─ web/         # React + Vite + TypeScript front-end
├─ docs/        # Generated Swagger sources
├─ docker-compose.yml
└─ Makefile     # DX commands (build, tests, docs, client gen)
```

- **Backend** exposes REST endpoints, enforces business rules, and publishes events for notifications.
- **Frontend** uses a generated client (`@appointment-master/api-client`) to consume the API safely.
- **Docs** are generated from Go annotations and used to regenerate the TypeScript client.
- **Docker Compose** runs Postgres, the API, and the web app together for parity with production.

---

## Appointment Types

- **Single**: slot-based with 1 attendee per slot — best for 1:1 sessions.
- **Group**: slot-based with multiple attendees per slot (capacity enforced) — ideal for classes or workshops.
- **Party**: a single shared slot for a single-day (or overnight) appointment window with a shared capacity — suited for events and open houses.


## Anti‑Scalping Controls

Asiko includes built‑in anti‑scalping rules to prevent duplicate or abusive bookings:\n
- **Standard**: blocks duplicate bookings by email.\n
- **Strict**: blocks duplicate bookings by email and device, and supports creator approval flows.\n
These settings help owners balance accessibility with control depending on the event type.

---

## Backend Highlights (Go + Gin)

- **Layered organization**: `api` (transport) → `services` (business rules) → `repository` (GORM data access).
- **Event bus**: a simple in-memory event bus publishes booking/appointment events; handlers are decoupled from business logic.
- **Status scheduling**: a unified scheduler advances appointment and booking statuses based on time.
- **Error taxonomy**: consistent error envelopes with machine codes, request IDs, and field errors.
- **Notification providers**: pluggable email provider abstraction; AhaSend is the default, with a noop provider for local/dev.
- **Testing**: unit and integration tests rely on mocks (`mockery`) and SQL mock where appropriate.

---

## Frontend Highlights (React + Vite)

- **Feature-first layout** (`web/src/features/*`) for clear domain boundaries (auth, appointments, bookings, analytics, notifications).
- **Typed API client** generated from OpenAPI; keeps UI and backend in lockstep.
- **Form validation** with React Hook Form + Zod for deterministic client-side checks.
- **Theming** with tokenized light/dark themes and reusable UI primitives.

---

## Setup

### Prerequisites

- Go 1.23+
- Node 20+ / npm
- Docker & Docker Compose

### Environment

Create `backend/.env` with the following (minimum):

```env
PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=appointmentdb
DB_SSLMODE=disable
JWT_SECRET_KEY=change-me

# Email provider selection
EMAIL_PROVIDER=ahasend

# AhaSend configuration
AHASEND_ACCOUNT_ID=...
AHASEND_API_KEY=...
AHASEND_FROM_EMAIL=hello@yourdomain.com
AHASEND_FROM_NAME=Asiko
AHASEND_ENABLED=true
```

Notes:
- Set `EMAIL_PROVIDER=noop` for local/dev without sending email.
- Appointment/booking time validation is enforced server-side; client UI blocks past dates/times.

### Local Development

```bash
make dev          # backend with Air hot reload
make web-dev      # Vite dev server at http://localhost:5173
make docker-dev   # postgres + api + web
```

Useful commands:

- `make docs-gen` – regenerate Swagger specs.
- `make client-gen` – regenerate the web API client.
- `make test` – run backend tests.
- `make mocks` – regenerate mocks.

Ports (default):

- API: `http://localhost:8890` (docker), `http://localhost:8080` (local)
- Web: `http://localhost:5173`
- Postgres: `localhost:5433`

---

## Software Practices Implemented

- **Contract-first API**: Swagger → generated TypeScript client.
- **Structured errors**: consistent API error envelope and logging.
- **Event-driven notifications**: services publish events; notification handlers subscribe.
- **Deterministic status transitions**: scheduler enforces time-based status updates.
- **Validation at both layers**: Zod (client) + server validation.
- **Test coverage**: repositories, services, and core business flows have unit tests.

---

## Design Choices

1. **Slot-based model** for single and group appointments for clear availability.
2. **Single-slot model** for parties with shared capacity to simplify booking rules.
3. **Event bus** keeps notifications decoupled from booking logic.
4. **Provider abstraction** for email allows swapping AhaSend without code changes.
5. **Generated client** reduces drift between frontend and backend.

---

## Project Structure (Detailed)

```
backend/
├─ api/                # Gin handlers & routing glue
├─ db/                 # Connection management + SQL migrations
├─ docs/               # swagger.go/json/yaml (generated)
├─ errors/             # AppError definitions, code maps, HTTP helpers
├─ middleware/         # request_id, logging, auth, error handler, CORS
├─ models/             # entities (GORM), requests, responses
├─ notifications/      # providers + templates + event handlers
├─ repository/         # GORM-backed persistence + mocks
├─ services/           # business logic + schedulers
├─ utils/              # helpers (codes, time math, validation)
└─ main.go             # wire-up (DB, repos, services, router)

web/
├─ api-client/         # Generated TypeScript client
├─ src/
│  ├─ app/             # router, providers, shell
│  ├─ components/      # shared UI primitives
│  ├─ features/        # domain slices (auth, appointments, bookings…)
│  ├─ services/        # API hooks, auth storage
│  └─ theme/           # design tokens + global styles
└─ scripts/            # client generation helper
```

---

## API Documentation

- Swagger UI is available at `/swagger/index.html` when the API is running.
- Specs live in `docs/swagger.{json,yaml}` and are generated via `make docs-gen`.

---
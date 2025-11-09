# Asiko Platform

Asiko is a full-stack booking platform that helps solo operators and small teams publish availability, collect reservations, and keep both sides of the marketplace in sync. The project recently evolved from a backend-only API into a deployable monorepo with shared tooling, a modern React front-end, and reproducible infrastructure so we can ship the hosted version with confidence.

---

## Why This Exists

Service businesses usually juggle ad-hoc spreadsheets, messaging threads, and siloed tools just to answer four questions: *What slots are open? Who booked them? Did we confirm? Can we trust the data?* Asiko tackles those pain points by:

- Keeping a single source of truth for appointments, bookings, ban lists, and analytics.
- Serving both authenticated business owners and lightweight guest flows in the same API.
- Providing a pre-built web experience that speaks to the API through a generated client, ensuring UI and backend stay in lockstep.
- Baking in observability (request IDs, structured errors) so we can debug production incidents once the platform is hosted.

---

## Scope & Goals

- **Scheduling & Ticketing**: Create single, group, and party-style appointments with smart slot generation and capacity enforcement.
- **Customer Access**: Users can browse via code, self-book as guests, or manage bookings when authenticated.
- **Operational Controls**: Ban lists, analytics, and notification hooks keep appointment owners in control.
- **Deployment Readiness**: Docker-first workflow, reproducible Make targets, and Swagger/OpenAPI-powered documentation.

---

## Architecture Overview

```
root
├─ backend/     # Go + Gin API (Clean Architecture layering)
├─ web/         # React + Vite + TypeScript front-end
├─ docs/        # Generated Swagger sources served statically
├─ docker-compose.yml
├─ Makefile     # Cross-cutting DX commands (build, tests, docs, client gen)
└─ env.*        # Shared environment maps for compose/containers
```

- **Backend** exposes REST endpoints, handles persistence, and raises domain-specific `AppError`s that map cleanly to HTTP responses.
- **Frontend** consumes the same contract via a generated `@appointment-master/api-client`, with feature-based routing, TanStack Query caching, and form validation via React Hook Form + Zod.
- **Docs** pipeline (`swag` + `openapi-generator`) keeps the schema authoritative, feeding both Swagger UI and the TypeScript client.
- **Docker Compose** runs Postgres, the Go API, and the Vite dev server together with hot reload support for hosting previews.

---

## Backend Highlights (Go 1.23 + Gin)

- **Layered organization**: `api` (transport) → `services` (business rules) → `repository` (GORM data access) keeps responsibilities isolated and testable.
- **Custom error taxonomy**: `errors/apperror.go` standardizes codes, kinds, and fields, surfaced through `middleware/error_handler.go`.
- **Middlewares**: request IDs, structured logging, CORS, auth, and panic-safe error handling wrap every handler.
- **Domain modules**: appointments, bookings, analytics, ban lists, pending users, and notification dispatch live in dedicated packages with mocks for unit tests.
- **Infrastructure**: `db/` hosts connection pooling and raw SQL migrations; `notifications/` abstracts providers (SendGrid today).
- **Testing**: service and handler suites rely on generated mocks (`mockery` via `make mocks`) to cover analytics, booking, and auth flows.

---

## Frontend Highlights (React 18 + Vite)

- **Feature-first layout** (`web/src/features/*`) keeps pages, hooks, and UI state close to their domains (auth, appointments, bookings, analytics, ban list, notifications).
- **App shell & routing**: `web/src/app/router.tsx` wires public, protected, and error routes; `ProtectedRoute` gates dashboards based on auth state.
- **State & forms**: TanStack Query for server cache, React Hook Form + Zod for validation, and React Hot Toast for UX feedback.
- **Styling**: `styled-components` with theme tokens in `web/src/theme`.
- **API fidelity**: `scripts/generate-api-client.sh` compiles the OpenAPI-driven client and TypeScript declarations so the UI can consume backend changes without drift.

---

## Cross-Cutting Decisions

1. **Monorepo DX** – A single `Makefile` proxies backend commands (dev server via Air, migrations, docs generation) and frontend dev scripts, so CI/CD can rely on consistent entry points.
2. **API-first development** – Swagger docs (`docs/swagger.json|yaml`) are generated from Go annotations, then used to regenerate the TypeScript client and publish docs.
3. **Strict error mapping** – Every layer returns `AppError` instances; middleware translates them into JSON envelopes with status, machine codes, request IDs, and field errors.
4. **Event notifications** – The `services/event_notification.go` abstraction routes booking events to the notification repository/service so we can add channels (email/SMS) later without touching handlers.
5. **Security defaults** – JWT auth for owner flows, ban lists for abuse mitigation, scoped request logging, and configurable rate-limiting hooks built into the error taxonomy.
6. **Docker as the contract** – `docker-compose.yml` is now the canonical way to run Postgres, API, and web together, mirroring the production topology we intend to host.

---

## Feature Set

- **Authentication & Verification**
  - Email-based signup/login with bcrypt hashing.
  - Pending-user verification flow with re-sent notifications.
  - JWT-protected endpoints and request-scoped auth middleware.
- **Appointment Management**
  - Single/group/party templates with automatic slot schedules.
  - Configurable capacity, duration, anti-scalping policies, and unique app codes.
  - Owner-centric listings (`/appointments/my`) and detail views.
- **Booking Lifecycle**
  - Guest and registered booking flows with slot availability checks.
  - Booking codes, status tracking, updates, cancellations, and analytics aggregation.
  - Ban list enforcement and device token checks to prevent abusive bookings.
- **Analytics & Operations**
  - Owner dashboards for bookings, revenue proxies, and slot utilization.
  - Ban list CRUD, notification preferences, and soon-to-land payment hooks.

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
├─ notifications/      # interfaces + SendGrid implementation
├─ repository/         # GORM-backed persistence + mocks
├─ services/           # Business logic + event notifications
├─ utils/              # helpers (codes, time math, validation)
└─ main.go             # wire-up (DB, repos, services, router, swagger)

web/
├─ api-client/         # Generated TypeScript client (OpenAPI Fetch)
├─ src/
│  ├─ app/             # router, providers, shell, layout styles
│  ├─ components/      # shared UI primitives
│  ├─ features/        # domain slices (auth, appointments, bookings…)
│  ├─ services/        # API hooks, query keys, auth storage
│  └─ theme/           # design tokens + global styles
└─ scripts/            # client generation helper
```

---

## Getting Started

### Prerequisites

- Go 1.23+
- Node 20+ / npm
- Docker & Docker Compose
- PostgreSQL client tools (optional but useful for debugging)

### Environment

Copy `backend/.env.example` (or create `backend/.env`) and set the following:

```env
PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=appointmentdb
DB_SSLMODE=disable
JWT_SECRET_KEY=change-me
SENDGRID_API_KEY=<optional>
```

`env.routes` exposes shared values to Docker services (API + Web). Local CLI commands will read from `backend/.env`.

### Local Development

```bash
make dev          # backend with Air hot reload
make web-dev      # Vite dev server at http://localhost:5173
make docker-dev   # postgres + api + web (mirrors hosting stack)
```

Useful extras:

- `make docs-gen` – refresh Swagger specs (`docs/swagger.json|yaml`).
- `make client-gen` – regenerate the TypeScript API client under `web/api-client`.
- `make test` – run Go unit tests with verbose output.
- `make mocks` – rebuild repository/service mocks via mockery.
- `make migrate-<up|down|create|to>` – manage SQL migrations.

Ports (default):

- API: `http://localhost:8890` when dockerized, `http://localhost:8080` when running via `make dev`.
- Web: `http://localhost:5173`
- Postgres: `localhost:5433`

---

## Deployment & Hosting Notes

- **Containers**: The production deployment mirrors `docker-compose.yml` (Postgres 17, Go API, React build served by Vite preview/static server). Each service sets `env_file` for secrets and exposes health checks (API `/health`, Postgres `pg_isready`).
- **Hot reload in Docker**: `develop.watch` entries rebuild the API and web images on file changes, which is useful in staging environments.
- **Static assets**: The web Dockerfile builds the Vite bundle; nginx or any static file server can serve `dist/` once deployed.
- **Migrations**: Run through `make migrate-up` (requires `DB_*` vars). CI/CD should execute migrations before bringing the API online.

---

## API Documentation

- Swagger UI is mounted at `/swagger/index.html` when the API is running.
- Source specs live in `docs/swagger.{json,yaml}` and are regenerated via `make docs-gen`.
- The front-end client pulls from those specs (`openapi-generator` → `web/api-client/src`), ensuring the UI type system matches backend contracts.

---

## Roadmap & Future Enhancements

- Payment integrations (Stripe/Paystack) for prepaid bookings.
- Notification fan-out beyond SendGrid (e.g., SMS, WhatsApp).
- Booking approval workflows with audit trails.
- Org/multi-tenant separation for agencies managing multiple brands.
- Expanded analytics (conversion funnels, attrition tracking).

---

## Contributing

1. Fork and branch from `main`.
2. Use the existing patterns (Clean Architecture, error taxonomy, feature-based React slices).
3. Add/refresh tests for any behavioral changes.
4. Regenerate docs/client if you modify the API surface.
5. Submit a PR referencing any relevant bug tracker entries.

---

## License

MIT © Asiko Contributors

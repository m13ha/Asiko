# Asiko Web (Vite + React + TypeScript)

This document is the development plan for the new web frontend that replaces the React Native app. It uses the generated API client in `api-client` for all backend access.

## Objectives
- Fast, typed web UI with Vite + React + TypeScript.
- First‑class DX: file‑based features, server state with TanStack Query, form validation with Zod.
- Strict use of the generated client `@appointment-master/api-client` for all API calls.
- Consistent theming, responsive layouts, and accessible components.

## Tech Stack
- App: Vite, React 18+, TypeScript
- Routing: React Router v6
- Data: TanStack Query v5
- Forms: React Hook Form + Zod
- UI/Theming: Styled-components (or Emotion) + CSS variables; react-hot-toast for toasts
- Dates: date-fns
- Icons: Lucide React
- Env: Vite envs (`VITE_API_BASE_URL`)

## Project Layout (web/)
- web/
  - index.html
  - public/
  - src/
    - app/
      - App.tsx (root shell)
      - router.tsx (route objects)
      - providers/
        - QueryProvider.tsx
        - ThemeProvider.tsx
      - styles/
        - globals.css (normalize + CSS vars)
    - features/
      - auth/ (login, signup, verify)
        - api.ts (wrap generated client calls)
        - hooks.ts (useLogin, useSignup, useVerify)
        - pages/ (LoginPage, SignupPage, VerifyPage)
        - components/ (AuthForm pieces)
      - appointments/
        - api.ts (create/list/details/users)
        - hooks.ts (useCreateAppointment, useMyAppointments, useAppointment, useUsersForAppointment)
        - pages/ (MyAppointmentsPage, CreateAppointmentPage, AppointmentDetailsPage)
        - components/ (AppointmentForm, AppointmentCard)
      - bookings/
        - api.ts (slots, book guest/registered, get/update/cancel by code)
        - hooks.ts (useAvailableSlots, useBook, useBookingByCode, useUpdateBooking, useCancelBooking)
        - pages/ (BookPage, BookingManagePage, MyBookingsPage)
        - components/ (BookingForm, SlotPicker, BookingSummary)
      - analytics/ (MVP+)
      - banlist/ (MVP+)
      - notifications/ (MVP+)
    - components/ (shared UI: Button, Input, Select, Card, Modal, Table, Tabs, Badge, Skeleton, Spinner, DatePicker, Toggle, Toast)
    - services/
      - auth.ts (token storage/helpers)
    - theme/
      - light.ts, dark.ts, index.ts (tokens)
    - types/, utils/
    - main.tsx

## App theme colors
  `#F9F7F7`
  `#DBE2EF`
  `#3F72AF`
  `#112D4E`

## Using the Generated API Client
We consume the client from `web/api-client` as a local dependency and configure it per environment.

- Install locally (from web/):
  - `npm i`
  - `npm i @appointment-master/api-client@file:./api-client`
  - Regenerate client from backend when API changes: `npm run client:gen`
- Configure the client:
  ```ts
  import { Configuration, AppointmentsApi, AuthenticationApi } from '@appointment-master/api-client';

  const tokenProvider = () => localStorage.getItem('token') || '';

  export const apiConfig = new Configuration({
    basePath: import.meta.env.VITE_API_BASE_URL,
    apiKey: () => `Bearer ${tokenProvider()}`, // Inject Authorization header
  });

  export const authApi = new AuthenticationApi(apiConfig);
  export const appointmentsApi = new AppointmentsApi(apiConfig);
  // ... other APIs from api-client/src/apis
  ```
- Pagination: Many endpoints return `{ items, page, per_page, total, total_pages }`.
  - Use types from the client (e.g., `GetMyAppointments200Response`).
  - You may map to UI-friendly shapes in feature `api.ts` wrappers.

## Environment & CORS
- Create `web/.env.local`:
  ```env
  VITE_API_BASE_URL=https://7obnqz8ix1.loclx.io
  ```
- Backend CORS for dev: set `CORS_ALLOWED_ORIGINS=http://localhost:5173,https://7obnqz8ix1.loclx.io`.

## Auth Strategy
- Store JWT in `localStorage`.
- Read token in `apiKey` callback for the client (above).
- Guard protected routes with a simple hook (`useAuth`) and wrapper.
- Optional later: add refresh flow via backend refresh endpoint.

## Development Phases & Deliverables

Phase 0 – Setup & Foundations (1–2 days)
- Scaffold Vite React TS app under `web/`.
- Install core deps (router, query, rhf, zod, styled-components, toast, date-fns, lucide-react).
- Add base providers (ThemeProvider, QueryClientProvider) and global styles.
- Integrate `@appointment-master/api-client` with runtime Configuration (env base URL, Authorization via apiKey).
- Result: App boots, health check page calls backend via generated client.

Phase 1 – Authentication (2–3 days)
- Pages: Login, Signup, Verify Registration.
- Hooks: `useLogin`, `useSignup`, `useVerify` (TanStack Query mutations).
- Persist token; redirect on auth; error/success toasts.
- Deliverables: Auth flows fully functional against backend.

Phase 2 – Appointments (Owner) (3–4 days)
- Pages: My Appointments (paginated), Create Appointment, Appointment Details.
- Components: AppointmentForm, AppointmentCard.
- Hooks for create/list/details, attendees/users per appointment.
- Deliverables: Owners can create and view appointments with `app_code` surfaced.

Phase 3 – Bookings (Public & Registered) (4–5 days)
- Public Booking: `/book` (enter code) → fetch appointment + slots → select date/slot → guest details → confirm; handle device token for strict anti‑scalping via `/auth/device-token`.
- My Bookings (registered): list `/appointments/registered`.
- Manage by Code: view/update/cancel via `/bookings/:booking_code`.
- Deliverables: Full booking lifecycle for guest and registered paths.

Phase 4 – Analytics, Ban List, Notifications (MVP+) (3–5 days)
- Analytics: `/analytics` overview; charts as needed.
- Ban List: add/remove/list.
- Notifications: list + mark read.
- Deliverables: Owner productivity features in place.

Phase 5 – Polish & Quality (ongoing)
- Theming (dark/light), responsive polish, a11y pass.
- 404/500 pages, error boundaries.
- Testing: unit tests for hooks/components; e2e smoke (Playwright) for core flows.
- Deliverables: Stable, maintainable UI ready for staging.

## Milestones & Acceptance Criteria
- M0: App boots; API client wired; `.env` works.
- M1: Auth flows stable; token persists; protected routes enforced.
- M2: Appointment CRUD (owner) complete; can view `app_code`.
- M3: Booking flows (guest/registered) complete; manage by `booking_code` works.
- M4: Analytics/Ban List/Notifications operational (MVP level).
- M5: Theming, a11y, tests; app production‑ready.

## Dev Scripts (suggested)
- `npm run dev` – Vite dev server
- `npm run build` – production build
- `npm run preview` – preview build
- `npm run typecheck` – TS check
- `npm run test` – unit tests (vitest/react-testing-library)

## Notes & Risks
- Ensure backend CORS allows Vite origin in dev, no wildcard fallback.
- Align time/date formats with backend expectations (ISO 8601).
- Use generated client models; avoid hand‑rolled fetch to keep in sync.
- If publishing the client, consider a workspace or GH registry; for now, use `file:` dependency.

---

When ready, I can scaffold `web/` with this structure, wire the providers, and add initial routes and API client configuration.

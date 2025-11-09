import { createBrowserRouter } from 'react-router-dom';
import { App } from './shell/App';
import { BookByCodePage } from '@/features/bookings/pages/BookByCodePage';
import { LoginPage } from '@/features/auth/pages/LoginPage';
import { SignupPage } from '@/features/auth/pages/SignupPage';
import { VerifyPage } from '@/features/auth/pages/VerifyPage';
import { ProtectedRoute } from '@/features/auth/ProtectedRoute';
import { DashboardPage } from '@/features/dashboard/pages/DashboardPage';
import { MyAppointmentsPage } from '@/features/appointments/pages/MyAppointmentsPage';
import { CreateAppointmentPage } from '@/features/appointments/pages/CreateAppointmentPage';
import { AppointmentDetailsPage } from '@/features/appointments/pages/AppointmentDetailsPage';
import { BookingManagePage } from '@/features/bookings/pages/BookingManagePage';
import { MyBookingsPage } from '@/features/bookings/pages/MyBookingsPage';
import { AnalyticsPage } from '@/features/analytics/pages/AnalyticsPage';
import { BanListPage } from '@/features/banlist/pages/BanListPage';
import { NotificationsPage } from '@/features/notifications/pages/NotificationsPage';
import { NotFoundPage } from '@/app/pages/NotFoundPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <BookByCodePage /> },
      { path: 'login', element: <LoginPage /> },
      { path: 'signup', element: <SignupPage /> },
      { path: 'verify', element: <VerifyPage /> },
      { path: 'bookings/:bookingCode', element: <BookingManagePage /> },
      { path: '*', element: <NotFoundPage /> },
      {
        element: <ProtectedRoute />,
        children: [
          { path: 'dashboard', element: <DashboardPage /> },
          { path: 'appointments', element: <MyAppointmentsPage /> },
          { path: 'appointments/new', element: <CreateAppointmentPage /> },
          { path: 'appointments/:id', element: <AppointmentDetailsPage /> },
          { path: 'bookings', element: <MyBookingsPage /> },
          { path: 'analytics', element: <AnalyticsPage /> },
          { path: 'ban-list', element: <BanListPage /> },
          { path: 'notifications', element: <NotificationsPage /> },
        ],
      },
      // Future: protected routes for dashboard, appointments, bookings
    ],
  },
]);

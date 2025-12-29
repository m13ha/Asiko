import { createBrowserRouter } from 'react-router-dom';
import { App } from './shell/App';
import { SplitPanelBookByCodePage } from '@/features/bookings/pages/SplitPanelBookByCodePage';
import { BookByCodePage } from '@/features/bookings/pages/BookByCodePage';
import { LoginPage } from '@/features/auth/pages/LoginPage';
import { SignupPage } from '@/features/auth/pages/SignupPage';
import { VerifyPage } from '@/features/auth/pages/VerifyPage';
import { ForgotPasswordPage } from '@/features/auth/pages/ForgotPasswordPage';
import { ResetPasswordPage } from '@/features/auth/pages/ResetPasswordPage';
import { ProtectedRoute } from '@/features/auth/ProtectedRoute';
import { HomeInsightsPage } from '@/features/dashboard/pages/HomeInsightsPage';
import { MyAppointmentsPage } from '@/features/appointments/pages/MyAppointmentsPage';
import { CreateAppointmentPage } from '@/features/appointments/pages/CreateAppointmentPage';
import { AppointmentDetailsPage } from '@/features/appointments/pages/AppointmentDetailsPage';
import { EditAppointmentPage } from '@/features/appointments/pages/EditAppointmentPage';
import { BookingManagePage } from '@/features/bookings/pages/BookingManagePage';
import { MyBookingsPage } from '@/features/bookings/pages/MyBookingsPage';
import { BanListPage } from '@/features/banlist/pages/BanListPage';
import { NotificationsPage } from '@/features/notifications/pages/NotificationsPage';
import { NotFoundPage } from '@/app/pages/NotFoundPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <SplitPanelBookByCodePage /> },
      { path: 'book-by-code', element: <SplitPanelBookByCodePage /> },
      { path: 'book-classic', element: <BookByCodePage /> },
      { path: 'login', element: <LoginPage /> },
      { path: 'signup', element: <SignupPage /> },
      { path: 'verify', element: <VerifyPage /> },
      { path: 'forgot-password', element: <ForgotPasswordPage /> },
      { path: 'reset-password', element: <ResetPasswordPage /> },
      { path: 'bookings/:bookingCode', element: <BookingManagePage /> },
      { path: '*', element: <NotFoundPage /> },
      {
        element: <ProtectedRoute />,
        children: [
          { path: 'dashboard', element: <HomeInsightsPage /> },
          { path: 'appointments', element: <MyAppointmentsPage /> },
          { path: 'appointments/new', element: <CreateAppointmentPage /> },
          { path: 'appointments/:id', element: <AppointmentDetailsPage /> },
          { path: 'appointments/:id/edit', element: <EditAppointmentPage /> },
          { path: 'bookings', element: <MyBookingsPage /> },
          { path: 'ban-list', element: <BanListPage /> },
          { path: 'notifications', element: <NotificationsPage /> },
        ],
      },
      // Future: protected routes for dashboard, appointments, bookings
    ],
  },
]);

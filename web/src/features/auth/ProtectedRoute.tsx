import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useIsAuthed } from '@/stores/authStore';

export function ProtectedRoute() {
  const isAuthed = useIsAuthed();
  const loc = useLocation();
  if (!isAuthed) return <Navigate to="/login" replace state={{ from: loc }} />;
  return <Outlet />;
}


import { useNavigate, Link } from 'react-router-dom';
import { useMyAppointments } from '../hooks';
import { Button } from '@/components/Button';
import { AppointmentCard } from '../components/AppointmentCard';
import { EmptyState, EmptyTitle, EmptyDescription, EmptyAction } from '@/components/EmptyState';

export function MyAppointmentsPage() {
  const { data, isLoading, error } = useMyAppointments();
  const navigate = useNavigate();

  return (
    <div style={{ display: 'grid', gap: 12 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0 }}>My Appointments</h1>
        <Button variant="primary" onClick={() => navigate('/appointments/new')}>Create Appointment</Button>
      </div>
      {isLoading && <div>Loading...</div>}
      {error && <div style={{ color: 'var(--danger)' }}>Failed to load appointments.</div>}
      <div style={{ display: 'grid', gap: 12, gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))' }}>
        {data?.items?.length ? (
          data.items.map((it: any) => <AppointmentCard key={it.id} item={it} />)
        ) : (
          <EmptyState>
            <EmptyTitle>No appointments yet</EmptyTitle>
            <EmptyDescription>Create your first appointment and share the code.</EmptyDescription>
            <EmptyAction>
              <Button variant="primary" onClick={() => navigate('/appointments/new')}>Create appointment</Button>
            </EmptyAction>
          </EmptyState>
        )}
      </div>
    </div>
  );
}

import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import * as API from '@appointment-master/api-client';
import { useMyAppointments } from '../hooks';
import { Button } from '@/components/Button';
import { AppointmentCard } from '../components/AppointmentCard';
import { EmptyState, EmptyTitle, EmptyDescription, EmptyAction } from '@/components/EmptyState';

const statusOptions = [
  { label: 'Pending', value: API.EntitiesAppointmentStatus.Pending },
  { label: 'Ongoing', value: API.EntitiesAppointmentStatus.Ongoing },
  { label: 'Completed', value: API.EntitiesAppointmentStatus.Completed },
  { label: 'Canceled', value: API.EntitiesAppointmentStatus.Canceled },
  { label: 'Expired', value: API.EntitiesAppointmentStatus.Expired },
] as const;

export function MyAppointmentsPage() {
  const [selectedStatuses, setSelectedStatuses] = useState<API.EntitiesAppointmentStatus[]>([]);
  const { data, isLoading, error } = useMyAppointments({ statuses: selectedStatuses });
  const navigate = useNavigate();
  const hasFilters = selectedStatuses.length > 0;

  const toggleStatus = (value: API.EntitiesAppointmentStatus) => {
    setSelectedStatuses((prev) =>
      prev.includes(value) ? prev.filter((status) => status !== value) : [...prev, value]
    );
  };

  return (
    <div style={{ display: 'grid', gap: 12 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0 }}>My Appointments</h1>
        <Button variant="primary" onClick={() => navigate('/appointments/new')}>Create Appointment</Button>
      </div>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, alignItems: 'center' }}>
        <span style={{ fontSize: 14, fontWeight: 600 }}>Status:</span>
        {statusOptions.map((option) => {
          const active = selectedStatuses.includes(option.value);
          return (
            <Button
              key={option.value}
              variant={active ? 'primary' : 'ghost'}
              onClick={() => toggleStatus(option.value)}
              style={{ padding: '4px 12px', fontSize: 12 }}
            >
              {option.label}
            </Button>
          );
        })}
        {hasFilters && (
          <Button variant="ghost" onClick={() => setSelectedStatuses([])} style={{ padding: '4px 12px', fontSize: 12 }}>
            Clear
          </Button>
        )}
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

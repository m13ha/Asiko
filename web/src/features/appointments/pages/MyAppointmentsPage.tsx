import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import * as API from '@appointment-master/api-client';
import { useMyAppointments } from '../hooks';
import { Button } from '@/components/Button';
import { AppointmentCard } from '../components/AppointmentCard';
import { EmptyState, EmptyTitle, EmptyDescription, EmptyAction } from '@/components/EmptyState';
import { PaginatedGrid } from '@/components/PaginatedGrid';
import { usePagination } from '@/hooks/usePagination';
import { Dropdown } from 'primereact/dropdown';

const statusOptions = [
  { label: 'Pending', value: API.EntitiesAppointmentStatus.Pending },
  { label: 'Ongoing', value: API.EntitiesAppointmentStatus.Ongoing },
  { label: 'Completed', value: API.EntitiesAppointmentStatus.Completed },
  { label: 'Canceled', value: API.EntitiesAppointmentStatus.Canceled },
  { label: 'Expired', value: API.EntitiesAppointmentStatus.Expired },
];

export function MyAppointmentsPage() {
  const [selectedStatuses, setSelectedStatuses] = useState<API.EntitiesAppointmentStatus[]>([]);
  const pagination = usePagination(1, 10);
  const { data, isLoading, error } = useMyAppointments({ 
    statuses: selectedStatuses,
    ...pagination.params 
  });
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8 bg-white rounded-2xl shadow-lg p-6">
          <div className="flex justify-between items-center mb-4">
            <h1 className="text-3xl font-bold text-gray-800">My Appointments</h1>
            <Button variant="primary" onClick={() => navigate('/appointments/new')}>
              Create Appointment
            </Button>
          </div>
          <div className="flex flex-wrap gap-3 items-center">
            <span className="text-sm font-semibold text-gray-600">Filter by Status:</span>
            <Dropdown 
              value={selectedStatuses[0] || null} 
              onChange={(e) => setSelectedStatuses(e.value ? [e.value] : [])} 
              options={statusOptions} 
              optionLabel="label" 
              showClear 
              placeholder="All statuses"
              className="min-w-48"
            />
          </div>
        </div>
        
        <PaginatedGrid
          data={data}
          isLoading={isLoading}
          error={error}
          onPageChange={pagination.updatePage}
          renderItem={(item: any) => <AppointmentCard key={item.id} item={item} />}
          emptyState={
            <div className="bg-white rounded-2xl shadow-lg p-12 text-center">
              <EmptyState>
                <EmptyTitle>No appointments yet</EmptyTitle>
                <EmptyDescription>Create your first appointment and share the code.</EmptyDescription>
                <EmptyAction>
                  <Button variant="primary" onClick={() => navigate('/appointments/new')}>
                    Create appointment
                  </Button>
                </EmptyAction>
              </EmptyState>
            </div>
          }
        />
      </div>
    </div>
  );
}
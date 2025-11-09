import { Card, CardHeader, CardTitle } from '@/components/Card';
import { AppointmentForm, AppointmentFormValues } from '../components/AppointmentForm';
import { useCreateAppointment } from '../hooks';
import { useNavigate } from 'react-router-dom';

export function CreateAppointmentPage() {
  const create = useCreateAppointment();
  const navigate = useNavigate();

  const onSubmit = (v: AppointmentFormValues) => {
    create.mutate(
      {
        title: v.title,
        description: v.description,
        type: v.type,
        bookingDuration: v.bookingDuration,
        startDate: v.startDate,
        endDate: v.endDate,
        startTime: v.startTime,
        endTime: v.endTime,
        maxAttendees: v.maxAttendees,
      },
      { onSuccess: (res: any) => navigate(`/appointments/${res?.id}`, { state: { appointment: res } }) }
    );
  };

  return (
    <div style={{ maxWidth: 720, margin: '0 auto' }}>
      <Card>
        <CardHeader>
          <CardTitle>Create Appointment</CardTitle>
        </CardHeader>
        <AppointmentForm onSubmit={onSubmit} pending={create.isPending} />
      </Card>
    </div>
  );
}


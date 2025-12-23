import { Card, CardHeader, CardTitle } from '@/components/Card';
import { AppointmentForm, AppointmentFormValues } from '../components/AppointmentForm';
import { useCreateAppointment } from '../hooks';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';

const toIsoDate = (date: string) => new Date(`${date}T00:00:00Z`).toISOString();

const normalizeTime = (time: string) => {
  if (!time.includes(':')) return `${time}:00`;
  const segments = time.split(':');
  if (segments.length === 2) return `${segments[0]}:${segments[1]}:00`;
  return `${segments[0]}:${segments[1]}:${segments[2] || '00'}`;
};

const toIsoTime = (date: string, time: string) => new Date(`${date}T${normalizeTime(time)}Z`).toISOString();

export function CreateAppointmentPage() {
  const create = useCreateAppointment();
  const navigate = useNavigate();

  const onSubmit = (v: AppointmentFormValues) => {
    const startDateIso = toIsoDate(v.startDate);
    const endDateIso = toIsoDate(v.endDate);
    const startTimeIso = toIsoTime(v.startDate, v.startTime);
    const endTimeIso = toIsoTime(v.startDate, v.endTime);

    create.mutate(
      {
        title: v.title,
        description: v.description,
        type: v.type,
        bookingDuration: v.bookingDuration,
        startDate: startDateIso,
        endDate: endDateIso,
        startTime: startTimeIso,
        endTime: endTimeIso,
        maxAttendees: v.maxAttendees ?? 1,
        antiScalpingLevel: v.antiScalpingLevel,
      },
      { onSuccess: (res: any) => navigate(`/appointments/${res?.id}`, { state: { appointment: res } }) }
    );
  };

  return (
    <div className="px-0 sm:px-6">
      <Card className="p-0 sm:p-6">
        <AppointmentForm onSubmit={onSubmit} pending={create.isPending} />
      </Card>
    </div>
  );
}

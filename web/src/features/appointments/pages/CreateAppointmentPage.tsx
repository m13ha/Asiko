import { Card, CardHeader, CardTitle } from '@/components/Card';
import { AppointmentForm, AppointmentFormValues } from '../components/AppointmentForm';
import { useCreateAppointment } from '../hooks';
import { useNavigate } from 'react-router-dom';

const formatOffset = (date: Date) => {
  const offsetMinutes = -date.getTimezoneOffset();
  const sign = offsetMinutes >= 0 ? '+' : '-';
  const abs = Math.abs(offsetMinutes);
  const hours = String(Math.floor(abs / 60)).padStart(2, '0');
  const minutes = String(abs % 60).padStart(2, '0');
  return `${sign}${hours}:${minutes}`;
};

const toIsoDate = (date: string) => {
  const base = new Date(`${date}T00:00:00`);
  return `${date}T00:00:00${formatOffset(base)}`;
};

const normalizeTime = (time: string) => {
  if (!time.includes(':')) return `${time}:00`;
  const segments = time.split(':');
  if (segments.length === 2) return `${segments[0]}:${segments[1]}:00`;
  return `${segments[0]}:${segments[1]}:${segments[2] || '00'}`;
};

const toIsoTime = (date: string, time: string) => {
  const normalized = normalizeTime(time);
  const base = new Date(`${date}T${normalized}`);
  return `${date}T${normalized}${formatOffset(base)}`;
};

export function CreateAppointmentPage() {
  const create = useCreateAppointment();
  const navigate = useNavigate();

  const onSubmit = (v: AppointmentFormValues) => {
    const startDateIso = toIsoDate(v.startDate);
    const endDateIso = toIsoDate(v.endDate);
    const startTimeIso = toIsoTime(v.startDate, v.startTime);
    const endTimeIso = toIsoTime(v.endDate, v.endTime);

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
      {
        onSuccess: () => {
          window.setTimeout(() => {
            navigate('/appointments');
          }, 400);
        },
      }
    );
  };

  return (
    <div className="px-0 sm:px-6">
      <Card className="p-0 sm:p-6">
        <AppointmentForm onSubmit={onSubmit} pending={create.isPending} submitLabel="Create appointment" />
      </Card>
    </div>
  );
}

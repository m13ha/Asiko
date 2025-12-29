import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { format } from 'date-fns';
import { useMemo } from 'react';
import * as API from '@appointment-master/api-client';
import { Card, CardHeader, CardTitle } from '@/components/Card';
import { Button } from '@/components/Button';
import { AppointmentForm, AppointmentFormValues } from '../components/AppointmentForm';
import { useAppointmentUsers, useMyAppointments, useUpdateAppointment } from '../hooks';

const toIsoDate = (date: string) => new Date(`${date}T00:00:00Z`).toISOString();

const normalizeTime = (time: string) => {
  if (!time.includes(':')) return `${time}:00`;
  const segments = time.split(':');
  if (segments.length === 2) return `${segments[0]}:${segments[1]}:00`;
  return `${segments[0]}:${segments[1]}:${segments[2] || '00'}`;
};

const toIsoTime = (date: string, time: string) => new Date(`${date}T${normalizeTime(time)}Z`).toISOString();

const toDateInput = (value?: string) => {
  if (!value) return '';
  const d = new Date(value);
  return Number.isNaN(d.getTime()) ? '' : format(d, 'yyyy-MM-dd');
};

const toTimeInput = (value?: string) => {
  if (!value) return '';
  const d = new Date(value);
  return Number.isNaN(d.getTime()) ? '' : format(d, 'HH:mm');
};

export function EditAppointmentPage() {
  const { id = '' } = useParams();
  const navigate = useNavigate();
  const location = useLocation() as { state?: { appointment?: any } };
  const { data: myAppointments, isLoading } = useMyAppointments({ page: 1, size: 100 });
  const appointment = location.state?.appointment || myAppointments?.items?.find((item: any) => item.id === id);
  const { data: users, isLoading: usersLoading } = useAppointmentUsers(appointment?.appCode ?? '', { page: 1, size: 1 }, {
    enabled: !!appointment?.appCode,
  });
  const update = useUpdateAppointment();

  const hasBookings = Boolean((users?.total ?? users?.items?.length ?? 0) > 0);

  const initialValues: AppointmentFormValues = useMemo(
    () => ({
      title: appointment?.title ?? '',
      description: appointment?.description ?? '',
      type: appointment?.type ?? API.EntitiesAppointmentType.Single,
      bookingDuration: appointment?.bookingDuration ?? 30,
      startDate: toDateInput(appointment?.startDate),
      endDate: toDateInput(appointment?.endDate),
      startTime: toTimeInput(appointment?.startTime),
      endTime: toTimeInput(appointment?.endTime),
      maxAttendees: appointment?.maxAttendees ?? undefined,
      antiScalpingLevel: appointment?.antiScalpingLevel ?? API.EntitiesAntiScalpingLevel.ScalpingStandard,
    }),
    [
      appointment?.bookingDuration,
      appointment?.description,
      appointment?.endDate,
      appointment?.endTime,
      appointment?.maxAttendees,
      appointment?.startDate,
      appointment?.startTime,
      appointment?.title,
      appointment?.type,
      appointment?.antiScalpingLevel,
    ]
  );

  if (!appointment && isLoading) {
    return <div className="py-10 text-center text-sm text-[var(--text-muted)]">Loading appointment...</div>;
  }

  if (!appointment) {
    return (
      <Card className="p-6">
        <CardHeader>
          <CardTitle>Appointment not found</CardTitle>
        </CardHeader>
        <div className="text-sm text-[var(--text-muted)]">We couldn&apos;t find this appointment in your list.</div>
        <div className="mt-4">
          <Button variant="outline" onClick={() => navigate('/appointments')}>Back to appointments</Button>
        </div>
      </Card>
    );
  }

  if (usersLoading) {
    return <div className="py-10 text-center text-sm text-[var(--text-muted)]">Loading bookings...</div>;
  }

  if (hasBookings) {
    return (
      <Card className="p-6">
        <CardHeader>
          <CardTitle>Editing locked</CardTitle>
        </CardHeader>
        <div className="text-sm text-[var(--text-muted)]">
          This appointment already has bookings. Updates are disabled, but you can still delete it from the details page.
        </div>
        <div className="mt-4 flex flex-wrap gap-2">
          <Button variant="outline" onClick={() => navigate(`/appointments/${appointment.id}`, { state: { appointment } })}>
            Back to appointment
          </Button>
        </div>
      </Card>
    );
  }

  const onSubmit = (v: AppointmentFormValues) => {
    // Temporary debug logs to verify form values and payload.
    console.log('[EditAppointment] raw form values', v);
    const merged = {
      title: v.title || appointment.title || '',
      description: v.description ?? appointment.description ?? '',
      type: v.type ?? appointment.type ?? API.EntitiesAppointmentType.Single,
      bookingDuration: v.bookingDuration ?? appointment.bookingDuration ?? 30,
      startDate: v.startDate || toDateInput(appointment.startDate),
      endDate: v.endDate || toDateInput(appointment.endDate),
      startTime: v.startTime || toTimeInput(appointment.startTime),
      endTime: v.endTime || toTimeInput(appointment.endTime),
      maxAttendees: v.maxAttendees ?? appointment.maxAttendees ?? 1,
      antiScalpingLevel: v.antiScalpingLevel ?? appointment.antiScalpingLevel ?? API.EntitiesAntiScalpingLevel.ScalpingStandard,
    };

    if (!merged.startDate || !merged.endDate || !merged.startTime || !merged.endTime) {
      console.log('[EditAppointment] missing schedule fields', merged);
      return;
    }

    const startDateIso = toIsoDate(merged.startDate);
    const endDateIso = toIsoDate(merged.endDate);
    const startTimeIso = toIsoTime(merged.startDate, merged.startTime);
    const endTimeIso = toIsoTime(merged.endDate, merged.endTime);

    console.log('[EditAppointment] payload', {
      title: merged.title,
      description: merged.description,
      type: merged.type,
      bookingDuration: merged.bookingDuration,
      startDate: startDateIso,
      endDate: endDateIso,
      startTime: startTimeIso,
      endTime: endTimeIso,
      maxAttendees: merged.maxAttendees,
      antiScalpingLevel: merged.antiScalpingLevel,
    });

    update.mutate(
      {
        id,
        input: {
          title: merged.title,
          description: merged.description,
          type: merged.type,
          bookingDuration: merged.bookingDuration,
          startDate: startDateIso,
          endDate: endDateIso,
          startTime: startTimeIso,
          endTime: endTimeIso,
          maxAttendees: merged.maxAttendees,
          antiScalpingLevel: merged.antiScalpingLevel,
        } as API.RequestsAppointmentRequest,
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
        <AppointmentForm
          onSubmit={onSubmit}
          pending={update.isPending}
          initialValues={initialValues}
          submitLabel="Save changes"
        />
      </Card>
    </div>
  );
}

import { useMutation, useQuery } from '@tanstack/react-query';
import * as API from '@appointment-master/api-client';
import toast from 'react-hot-toast';
import { createAppointment, getUsersForAppointment, listMyAppointments } from './api';

async function parseError(e: unknown): Promise<string> {
  if (e instanceof API.ResponseError) {
    try { const body = await e.response.json(); return body?.message || e.response.statusText || 'Request failed'; }
    catch { return e.response.statusText || 'Request failed'; }
  }
  return 'Something went wrong';
}

export function useMyAppointments(filters?: { statuses?: API.EntitiesAppointmentStatus[] }) {
  const statuses = filters?.statuses ?? [];
  const key = statuses.length ? statuses.slice().sort().join(',') : 'all';
  return useQuery({
    queryKey: ['my-appointments', key],
    queryFn: () => listMyAppointments({ statuses }),
  });
}

export function useCreateAppointment() {
  return useMutation({
    mutationFn: createAppointment,
    onSuccess: () => toast.success('Appointment created'),
    onError: async (e) => toast.error(await parseError(e)),
  });
}

export function useAppointmentUsers(id: string) {
  return useQuery({ queryKey: ['appointment-users', id], queryFn: () => getUsersForAppointment(id), enabled: !!id });
}

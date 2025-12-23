import { useMutation, useQuery } from '@tanstack/react-query';
import * as API from '@appointment-master/api-client';
import toast from 'react-hot-toast';
import { createAppointment, getAppointmentByAppCode, getUsersForAppointment, listMyAppointments } from './api';

async function parseError(e: unknown): Promise<string> {
  if (e instanceof API.ResponseError) {
    try { const body = await e.response.json(); return body?.message || e.response.statusText || 'Request failed'; }
    catch { return e.response.statusText || 'Request failed'; }
  }
  return 'Something went wrong';
}

export function useMyAppointments(filters?: {
  statuses?: API.EntitiesAppointmentStatus[];
  page?: number;
  size?: number;
}) {
  const statuses = filters?.statuses ?? [];
  const page = filters?.page ?? 0;
  const size = filters?.size ?? 10;
  const key = `${statuses.length ? statuses.slice().sort().join(',') : 'all'}-${page}-${size}`;
  return useQuery({
    queryKey: ['my-appointments', key],
    queryFn: () => listMyAppointments({ statuses, page, size }),
  });
}

export function useCreateAppointment() {
  return useMutation({
    mutationFn: createAppointment,
    onSuccess: () => toast.success('Appointment created'),
    onError: async (e) => toast.error(await parseError(e)),
  });
}

export function useAppointmentUsers(id: string, params?: { page?: number; size?: number }, options?: { enabled?: boolean }) {
  const page = params?.page ?? 1;
  const size = params?.size ?? 10;
  return useQuery({
    queryKey: ['appointment-users', id, page, size],
    queryFn: () => getUsersForAppointment(id, { page, size }),
    enabled: !!id && (options?.enabled ?? true),
  });
}

export function useAppointmentByAppCode(appCode: string) {
  return useQuery({
    queryKey: ['appointment', appCode],
    queryFn: () => getAppointmentByAppCode(appCode),
    enabled: !!appCode,
  });
}

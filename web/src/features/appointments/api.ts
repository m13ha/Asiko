import type { AppointmentsApi, GetMyAppointments200Response, RequestsAppointmentRequest } from '@appointment-master/api-client';
import { appointmentsApi } from '@/services/api';

export type AppointmentType = string;

export type CreateAppointmentInput = RequestsAppointmentRequest;

export const appointmentsClient: AppointmentsApi = appointmentsApi as unknown as AppointmentsApi;

export function listMyAppointments() {
  return appointmentsClient.getMyAppointments();
}

export function createAppointment(input: CreateAppointmentInput) {
  return appointmentsClient.createAppointment({ appointment: input });
}

export function getUsersForAppointment(appCode: string) {
  return appointmentsClient.getUsersRegisteredForAppointment({ appCode });
}

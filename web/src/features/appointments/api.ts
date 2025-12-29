import type {
  AppointmentsApi,
  EntitiesAppointment,
  EntitiesAppointmentStatus,
  GetMyAppointments200Response,
  RequestsAppointmentRequest,
} from '@appointment-master/api-client';
import { EntitiesAppointmentFromJSON } from '@appointment-master/api-client';
import { appointmentsApi } from '@/services/api';

export type AppointmentType = string;

export type CreateAppointmentInput = RequestsAppointmentRequest;

export const appointmentsClient: AppointmentsApi = appointmentsApi as unknown as AppointmentsApi;

export function listMyAppointments(params?: {
  statuses?: EntitiesAppointmentStatus[];
  page?: number;
  size?: number;
}): Promise<GetMyAppointments200Response> {
  const request: any = {};
  if (params?.statuses?.length) request.status = params.statuses;
  if (params?.page) request.page = params.page;
  if (params?.size) request.size = params.size;
  return appointmentsClient.getMyAppointments(Object.keys(request).length ? request : undefined);
}

export function createAppointment(input: CreateAppointmentInput) {
  return appointmentsClient.createAppointment({ appointment: input });
}

export function getUsersForAppointment(appCode: string, params?: { page?: number; size?: number }) {
  const request: any = { appCode };
  if (params?.page) request.page = params.page;
  if (params?.size) request.size = params.size;
  return appointmentsClient.getUsersRegisteredForAppointment(request);
}

export function getAppointmentByAppCode(appCode: string): Promise<EntitiesAppointment> {
  return appointmentsClient.getAppointmentByAppCode({ appCode });
}

export function updateAppointment(id: string, input: RequestsAppointmentRequest): Promise<EntitiesAppointment> {
  return appointmentsClient.updateAppointment({
    id,
    appointment: input,
  });
}

export function deleteAppointment(id: string): Promise<EntitiesAppointment> {
  return appointmentsClient.deleteAppointment({ id });
}

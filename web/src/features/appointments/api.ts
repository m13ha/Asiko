import type {
  AppointmentsApi,
  EntitiesAppointment,
  EntitiesAppointmentStatus,
  GetMyAppointments200Response,
  RequestsAppointmentRequest,
} from '@appointment-master/api-client';
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

export function getUsersForAppointment(appCode: string) {
  return appointmentsClient.getUsersRegisteredForAppointment({ appCode });
}

export function getAppointmentByAppCode(appCode: string): Promise<EntitiesAppointment> {
  return appointmentsClient.getAppointmentByAppCode({ appCode });
}

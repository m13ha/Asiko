import * as API from '@appointment-master/api-client';
import { getToken } from './auth';

const basePath =
  import.meta.env.VITE_API_BASE_URL ?? "https://jrjik-102-88-113-9.a.free.pinggy.link";

export const apiConfig = new API.Configuration({
  basePath,
  apiKey: () => `Bearer ${getToken()}`,
});

export const authApi = new API.AuthenticationApi(apiConfig);
export const appointmentsApi = new API.AppointmentsApi(apiConfig);
export const bookingsApi = new API.BookingsApi(apiConfig);
export const analyticsApi = new API.AnalyticsApi(apiConfig);
export const notificationsApi = new API.NotificationsApi(apiConfig);
export const banListApi = new API.BanListApi(apiConfig);

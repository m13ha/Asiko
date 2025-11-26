import { notificationsApi } from '@/services/api';

export type NotificationItem = any;
export type Paginated<T> = { items?: T[]; page?: number; per_page?: number; total?: number; total_pages?: number };

export async function listNotifications(params?: { page?: number; size?: number }): Promise<Paginated<NotificationItem>> {
  return notificationsApi.getNotifications(params);
}

export async function markAllNotificationsRead(): Promise<{ message?: string }> {
  return notificationsApi.markAllNotificationsAsRead();
}

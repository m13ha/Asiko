import { notificationsApi } from '@/services/api';

export type NotificationItem = any;
export type Paginated<T> = { items?: T[]; page?: number; per_page?: number; total?: number; total_pages?: number };

export async function listNotifications(): Promise<Paginated<NotificationItem>> {
  return notificationsApi.getNotifications();
}

export async function markAllNotificationsRead(): Promise<{ message?: string }> {
  return notificationsApi.markAllNotificationsAsRead();
}

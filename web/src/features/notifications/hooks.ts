import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { listNotifications, markAllNotificationsRead } from './api';
import toast from 'react-hot-toast';

export function useNotifications() {
  return useQuery({ queryKey: ['notifications'], queryFn: listNotifications });
}

export function useMarkAllRead() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: markAllNotificationsRead,
    onSuccess: () => {
      toast.success('All notifications marked as read');
      qc.invalidateQueries({ queryKey: ['notifications'] });
    },
    onError: () => toast.error('Failed to mark as read'),
  });
}


import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { listNotifications, markAllNotificationsRead } from './api';
import toast from 'react-hot-toast';

export function useNotifications(params?: { page?: number; size?: number }) {
  const page = params?.page ?? 1;
  const size = params?.size ?? 10;
  return useQuery({ 
    queryKey: ['notifications', page, size], 
    queryFn: () => listNotifications({ page, size }) 
  });
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


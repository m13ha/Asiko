import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { addToBanList, getBanList, removeFromBanList } from './api';

export function useBanList() {
  return useQuery({ queryKey: ['ban-list'], queryFn: () => getBanList() });
}

export function useAddBan() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (email: string) => addToBanList(email),
    onSuccess: () => {
      toast.success('Email added to ban list');
      qc.invalidateQueries({ queryKey: ['ban-list'] });
    },
    onError: () => toast.error('Failed to add email'),
  });
}

export function useRemoveBan() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (email: string) => removeFromBanList(email),
    onSuccess: () => {
      toast.success('Email removed');
      qc.invalidateQueries({ queryKey: ['ban-list'] });
    },
    onError: () => toast.error('Failed to remove email'),
  });
}


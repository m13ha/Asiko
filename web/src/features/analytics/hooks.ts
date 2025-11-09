import { useQuery } from '@tanstack/react-query';
import * as API from '@appointment-master/api-client';
import toast from 'react-hot-toast';
import { getUserAnalytics } from './api';

export function useUserAnalytics(startDate: string, endDate: string) {
  return useQuery({
    queryKey: ['analytics', startDate, endDate],
    queryFn: async () => {
      try { return await getUserAnalytics({ startDate, endDate }); }
      catch (e) {
        if (e instanceof API.ResponseError) {
          try { const body = await e.response.json(); toast.error(body?.message || e.response.statusText); } catch { toast.error(e.response.statusText); }
        } else { toast.error('Failed to load analytics'); }
        throw e;
      }
    },
    enabled: !!startDate && !!endDate,
  });
}

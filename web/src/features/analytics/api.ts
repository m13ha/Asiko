import { analyticsApi } from '@/services/api';

export function getUserAnalytics(params: { startDate: string; endDate: string }) {
  return analyticsApi.getUserAnalytics(params);
}

export function getDashboardAnalytics(params: { startDate: string; endDate: string }) {
  return analyticsApi.getDashboardAnalytics(params);
}


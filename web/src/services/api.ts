import * as API from '@appointment-master/api-client';
import { useAuthStore, refreshTokensAction } from '@/stores/authStore';

const basePath =
  import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8890';

function decodeJwtExp(token: string): number | undefined {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) return undefined;
    const payload = JSON.parse(atob(parts[1]));
    if (payload && typeof payload.exp === 'number') {
      return payload.exp * 1000;
    }
  } catch {
    return undefined;
  }
}

function computeExpiresAt(token: string, expiresIn?: number): number | undefined {
  if (expiresIn && expiresIn > 0) return Date.now() + expiresIn * 1000;
  return decodeJwtExp(token);
}

const authMiddleware: API.Middleware = {
  pre: async () => {
    const { accessToken, refreshToken, expiresAt } = useAuthStore.getState();
    if (!refreshToken) return;
    
    // Proactive refresh: if token is missing or expires in < 60s
    if (!accessToken || (expiresAt && expiresAt - Date.now() < 60_000)) {
      try {
        await refreshTokensAction(async (rt) => {
          const refreshed = await authApi.refreshToken({
            refresh: { refreshToken: rt },
          });
          if (!refreshed.token) throw new Error('Missing token');
          return {
            accessToken: refreshed.token,
            refreshToken: refreshed.refreshToken,
            expiresAt: computeExpiresAt(refreshed.token, refreshed.expiresIn),
          };
        });
      } catch {
        useAuthStore.getState().clearTokens();
      }
    }
  },
  post: async (context) => {
    const status = context.response.status;
    const { refreshToken } = useAuthStore.getState();

    // Reactive refresh: 401/403 with a refresh token available
    if ((status === 401 || status === 403) && refreshToken) {
      const alreadyRetried = (context.init.headers as Record<string, string> | undefined)?.[
        'x-am-refresh-attempt'
      ];
      if (alreadyRetried) return context.response;

      try {
        const newToken = await refreshTokensAction(async (rt) => {
          const refreshed = await authApi.refreshToken({
            refresh: { refreshToken: rt },
          });
          if (!refreshed.token) throw new Error('Missing token');
          return {
            accessToken: refreshed.token,
            refreshToken: refreshed.refreshToken,
            expiresAt: computeExpiresAt(refreshed.token, refreshed.expiresIn),
          };
        });

        if (!newToken) throw new Error('Refresh failed');

        const headers: Record<string, string> = {
          ...(context.init.headers as Record<string, string> | undefined),
          Authorization: `Bearer ${newToken}`,
          'x-am-refresh-attempt': 'true',
        };
        return await context.fetch(context.url, { ...context.init, headers });
      } catch {
        useAuthStore.getState().clearTokens();
        return context.response;
      }
    }
    return context.response;
  },
};

export const apiConfig = new API.Configuration({
  basePath,
  apiKey: () => {
    const { accessToken } = useAuthStore.getState();
    return accessToken ? `Bearer ${accessToken}` : '';
  },
  middleware: [authMiddleware],
});

class RawApi extends API.BaseAPI {
  async jsonRequest<T>(opts: API.RequestOpts, mapper?: (payload: any) => T): Promise<T> {
    const authHeader = await buildAuthHeader();
    const response = await this.request({
      ...opts,
      headers: {
        ...authHeader,
        ...opts.headers,
      },
    });
    const payload = await response.json();
    return mapper ? mapper(payload) : (payload as T);
  }
}

async function buildAuthHeader(): Promise<API.HTTPHeaders> {
  if (!apiConfig.apiKey) return {};
  const token = await apiConfig.apiKey('Authorization');
  if (!token) return {};
  return { Authorization: token };
}

export const authApi = new API.AuthenticationApi(apiConfig);
export const appointmentsApi = new API.AppointmentsApi(apiConfig);
export const bookingsApi = new API.BookingsApi(apiConfig);
export const analyticsApi = new API.AnalyticsApi(apiConfig);
export const notificationsApi = new API.NotificationsApi(apiConfig);
export const banListApi = new API.BanListApi(apiConfig);
export const rawApi = new RawApi(apiConfig);

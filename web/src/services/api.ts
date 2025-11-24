import * as API from '@appointment-master/api-client';
import {
  clearTokens,
  getRefreshToken,
  getTokens,
  setTokens,
} from './auth';

const basePath =
  import.meta.env.VITE_API_BASE_URL ?? 'https://jrjik-102-88-113-9.a.free.pinggy.link';

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

let refreshPromise: Promise<string | void> | null = null;

async function refreshAccessToken(): Promise<string> {
  if (refreshPromise) {
    const token = await refreshPromise;
    return (token as string) || getTokens().accessToken;
  }

  const refreshToken = getRefreshToken();
  if (!refreshToken) {
    clearTokens();
    throw new Error('Missing refresh token');
  }

  refreshPromise = (async () => {
    const refreshed = await authApi.refreshToken({
      refresh: { refreshToken },
    });
    if (!refreshed.token) {
      clearTokens();
      throw new Error('Missing token');
    }
    const expiresAt = computeExpiresAt(refreshed.token, refreshed.expiresIn);
    setTokens({
      accessToken: refreshed.token,
      refreshToken: refreshed.refreshToken,
      expiresAt,
    });
    return refreshed.token;
  })();

  try {
    const token = await refreshPromise;
    return (token as string) || getTokens().accessToken;
  } finally {
    refreshPromise = null;
  }
}

const authMiddleware: API.Middleware = {
  pre: async () => {
    const { accessToken, refreshToken, expiresAt } = getTokens();
    if (!refreshToken) return;
    if (!accessToken || (expiresAt && expiresAt - Date.now() < 60_000)) {
      try {
        await refreshAccessToken();
      } catch {
        clearTokens();
      }
    }
  },
  post: async (context) => {
    const status = context.response.status;
    if ((status === 401 || status === 403) && getRefreshToken()) {
      const alreadyRetried = (context.init.headers as Record<string, string> | undefined)?.[
        'x-am-refresh-attempt'
      ];
      if (alreadyRetried) return context.response;

      try {
        await refreshAccessToken();
        const tokens = getTokens();
        const headers: Record<string, string> = {
          ...(context.init.headers as Record<string, string> | undefined),
          Authorization: tokens.accessToken ? `Bearer ${tokens.accessToken}` : '',
          'x-am-refresh-attempt': 'true',
        };
        return await context.fetch(context.url, { ...context.init, headers });
      } catch {
        clearTokens();
        return context.response;
      }
    }
    return context.response;
  },
};

export const apiConfig = new API.Configuration({
  basePath,
  apiKey: () => {
    const { accessToken } = getTokens();
    return accessToken ? `Bearer ${accessToken}` : '';
  },
  middleware: [authMiddleware],
});

export const authApi = new API.AuthenticationApi(apiConfig);
export const appointmentsApi = new API.AppointmentsApi(apiConfig);
export const bookingsApi = new API.BookingsApi(apiConfig);
export const analyticsApi = new API.AnalyticsApi(apiConfig);
export const notificationsApi = new API.NotificationsApi(apiConfig);
export const banListApi = new API.BanListApi(apiConfig);

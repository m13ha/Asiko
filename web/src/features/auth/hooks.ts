import { useMutation } from '@tanstack/react-query';
import { authApi } from '@/services/api';
import * as API from '@appointment-master/api-client';
import toast from 'react-hot-toast';
import { useAuth } from './AuthProvider';

function computeExpiresAt(token?: string, expiresIn?: number) {
  if (expiresIn && expiresIn > 0) return Date.now() + expiresIn * 1000;
  if (!token) return undefined;
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

async function parseError(e: unknown): Promise<string> {
  if (e instanceof API.ResponseError) {
    try {
      const body = await e.response.json();
      return body?.message || e.response.statusText || 'Request failed';
    } catch {
      return e.response.statusText || 'Request failed';
    }
  }
  return 'Something went wrong';
}

export function useLogin() {
  const { setTokens } = useAuth();
  return useMutation({
    mutationFn: async (vars: { email: string; password: string }) =>
      authApi.loginUser({ login: { email: vars.email, password: vars.password } }),
    onSuccess: (res: any) => {
      if (res?.token) {
        setTokens({
          accessToken: res.token,
          refreshToken: res.refreshToken,
          expiresAt: computeExpiresAt(res.token, res.expiresIn),
        });
      }
      toast.success('Logged in');
    },
    onError: async (e) => toast.error(await parseError(e)),
  });
}

export function useLogout() {
  const { logout } = useAuth();
  return useMutation({
    mutationFn: async () => authApi.logoutUser(),
    onSettled: () => {
      logout();
      toast.success('Logged out');
    },
  });
}

export function useSignup() {
  return useMutation({
    mutationFn: async (vars: { name: string; email: string; password: string }) =>
      authApi.createUser({ user: { name: vars.name, email: vars.email, password: vars.password } }),
    onSuccess: (res) => toast.success(res?.message ?? 'Registration pending. Check your email.'),
    onError: async (e) => toast.error(await parseError(e)),
  });
}

export function useVerify() {
  const { setTokens } = useAuth();
  return useMutation({
    mutationFn: async (vars: { email: string; code: string }) =>
      authApi.verifyRegistration({ verification: { email: vars.email, code: vars.code } }),
    onSuccess: (res) => {
      if (res?.token) {
        setTokens({
          accessToken: res.token,
          refreshToken: res.refreshToken,
          expiresAt: computeExpiresAt(res.token, res.expiresIn),
        });
      }
      toast.success('Email verified');
    },
    onError: async (e) => toast.error(await parseError(e)),
  });
}

export function useResendVerification() {
  return useMutation({
    mutationFn: async (vars: { email: string }) =>
      authApi.resendVerification({ resend: { email: vars.email } }),
    onSuccess: (res) => toast.success(res?.message ?? 'Verification code resent'),
    onError: async (e) => toast.error(await parseError(e)),
  });
}

export function useDeviceToken() {
  return useMutation({
    mutationFn: async (vars: { deviceId: string }) =>
      authApi.generateDeviceToken({ device: { deviceId: vars.deviceId } }),
  });
}

import { PropsWithChildren, createContext, useContext, useEffect, useMemo, useState } from 'react';
import {
  getTokens as readTokens,
  clearTokens as dropTokens,
  setTokens as writeTokens,
  type AuthTokens,
} from '@/services/auth';

type AuthCtx = {
  token: string;
  isAuthed: boolean;
  setToken: (t: string) => void;
  setTokens: (tokens: AuthTokens) => void;
  logout: () => void;
};

const Ctx = createContext<AuthCtx | undefined>(undefined);

export function AuthProvider({ children }: PropsWithChildren) {
  const [token, setTokenState] = useState<string>(() => readTokens().accessToken);

  useEffect(() => {
    const onStorage = (e: StorageEvent) => {
      if (e.key === 'token' || e.key === 'refreshToken' || e.key === 'tokenExpiresAt') {
        setTokenState(readTokens().accessToken);
      }
    };
    window.addEventListener('storage', onStorage);
    return () => window.removeEventListener('storage', onStorage);
  }, []);

  const setToken = (t: string) => {
    setTokens({ accessToken: t });
  };

  const setTokens = (tokens: AuthTokens) => {
    writeTokens(tokens);
    setTokenState(tokens.accessToken);
  };

  const logout = () => {
    dropTokens();
    setTokenState('');
  };

  const value = useMemo<AuthCtx>(
    () => ({ token, isAuthed: !!token, setToken, setTokens, logout }),
    [token]
  );

  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export function useAuth() {
  const ctx = useContext(Ctx);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}

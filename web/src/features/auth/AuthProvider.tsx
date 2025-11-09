import { PropsWithChildren, createContext, useContext, useEffect, useMemo, useState } from 'react';
import { getToken as readToken, clearToken as dropToken, setToken as writeToken } from '@/services/auth';

type AuthCtx = {
  token: string;
  isAuthed: boolean;
  setToken: (t: string) => void;
  logout: () => void;
};

const Ctx = createContext<AuthCtx | undefined>(undefined);

export function AuthProvider({ children }: PropsWithChildren) {
  const [token, setTokenState] = useState<string>(() => readToken());

  useEffect(() => {
    const onStorage = (e: StorageEvent) => {
      if (e.key === 'token') setTokenState(readToken());
    };
    window.addEventListener('storage', onStorage);
    return () => window.removeEventListener('storage', onStorage);
  }, []);

  const setToken = (t: string) => {
    writeToken(t);
    setTokenState(t);
  };

  const logout = () => {
    dropToken();
    setTokenState('');
  };

  const value = useMemo<AuthCtx>(() => ({ token, isAuthed: !!token, setToken, logout }), [token]);

  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export function useAuth() {
  const ctx = useContext(Ctx);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}


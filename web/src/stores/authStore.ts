import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

export type AuthTokens = {
  accessToken: string;
  refreshToken?: string;
  expiresAt?: number;
};

export type AuthStatus = 'loading' | 'authenticated' | 'unauthenticated';

interface AuthState {
  accessToken: string;
  refreshToken: string;
  expiresAt: number | null;
  status: AuthStatus;
  
  // Actions
  setTokens: (tokens: AuthTokens) => void;
  clearTokens: () => void;
  setStatus: (status: AuthStatus) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      accessToken: '',
      refreshToken: '',
      expiresAt: null,
      status: 'loading',

      setTokens: (tokens) => set({
        accessToken: tokens.accessToken,
        refreshToken: tokens.refreshToken ?? '',
        expiresAt: tokens.expiresAt ?? null,
        status: tokens.accessToken ? 'authenticated' : 'unauthenticated',
      }),

      clearTokens: () => set({
        accessToken: '',
        refreshToken: '',
        expiresAt: null,
        status: 'unauthenticated',
      }),

      setStatus: (status) => set({ status }),
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
      // Only persist tokens and expiry, status can be re-evaluated on hydration
      partialize: (state) => ({
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        expiresAt: state.expiresAt,
      }),
      onRehydrateStorage: () => (state) => {
        if (state) {
          state.status = state.accessToken ? 'authenticated' : 'unauthenticated';
        }
      },
    }
  )
);

// Derived selector
export const useIsAuthed = () => useAuthStore((state) => state.status === 'authenticated');

// Cross-tab synchronization
let isSyncing = false;
export function startCrossTabSync() {
  if (isSyncing) return;
  isSyncing = true;
  
  window.addEventListener('storage', (event) => {
    if (event.key === 'auth-storage') {
      // Re-hydrate the store from localStorage manually
      useAuthStore.persist.rehydrate();
    }
  });
}

// Single-flight refresh logic (ported from services/api.ts logic)
let refreshPromise: Promise<string | null> | null = null;

export async function refreshTokensAction(refreshFn: (refreshToken: string) => Promise<AuthTokens>) {
  if (refreshPromise) return refreshPromise;

  const { refreshToken } = useAuthStore.getState();
  if (!refreshToken) {
    useAuthStore.getState().clearTokens();
    return null;
  }

  refreshPromise = (async () => {
    try {
      const tokens = await refreshFn(refreshToken);
      useAuthStore.getState().setTokens(tokens);
      return tokens.accessToken;
    } catch (error) {
      useAuthStore.getState().clearTokens();
      return null;
    } finally {
      refreshPromise = null;
    }
  })();

  return refreshPromise;
}

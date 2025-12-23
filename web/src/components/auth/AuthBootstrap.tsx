import { useEffect, useRef } from 'react';
import { useAuthStore, startCrossTabSync } from '@/stores/authStore';

export function AuthBootstrap({ children }: { children: React.ReactNode }) {
  const initialized = useRef(false);

  useEffect(() => {
    if (initialized.current) return;
    initialized.current = true;

    // Start cross-tab synchronization
    startCrossTabSync();

    // Re-check status based on persisted tokens
    const { accessToken } = useAuthStore.getState();
    if (accessToken) {
      useAuthStore.getState().setStatus('authenticated');
    } else {
      useAuthStore.getState().setStatus('unauthenticated');
    }
  }, []);

  return <>{children}</>;
}

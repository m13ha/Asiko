import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';
import { QueryProvider } from '@/app/providers/QueryProvider';
import { ThemeProvider } from '@/app/providers/ThemeProvider';
import { router } from '@/app/router';
import '@/app/styles/globals.css';
import '@/app/styles/primereact-overrides.css';
import 'primereact/resources/themes/lara-light-blue/theme.css';
import 'primereact/resources/primereact.min.css';
import 'primeicons/primeicons.css';
import { AuthProvider } from '@/features/auth/AuthProvider';
import { ErrorBoundary } from '@/app/errors/ErrorBoundary';
import { PrimeReactProvider } from 'primereact/api';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
    <React.StrictMode>
      <PrimeReactProvider value={{ ripple: true }}>
        <ThemeProvider>
          <QueryProvider>
            <AuthProvider>
              <ErrorBoundary>
                <RouterProvider router={router} />
              </ErrorBoundary>
            </AuthProvider>
          </QueryProvider>
        </ThemeProvider>
      </PrimeReactProvider>
    </React.StrictMode>
);

import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';
import { QueryProvider } from '@/app/providers/QueryProvider';
import { ThemeProvider } from '@/app/providers/ThemeProvider';
import { router } from '@/app/router';
import '@/app/styles/globals.css';
import 'primereact/resources/themes/lara-light-blue/theme.css';
import 'primereact/resources/primereact.min.css';
import 'primeicons/primeicons.css';
import { ErrorBoundary } from '@/app/errors/ErrorBoundary';
import { PrimeReactProvider } from 'primereact/api';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

import { AuthBootstrap } from '@/components/auth/AuthBootstrap';

root.render(
    <React.StrictMode>
      <PrimeReactProvider value={{ ripple: true }}>
        <ThemeProvider>
          <QueryProvider>
            <AuthBootstrap>
              <ErrorBoundary>
                <RouterProvider router={router} />
              </ErrorBoundary>
            </AuthBootstrap>
          </QueryProvider>
        </ThemeProvider>
      </PrimeReactProvider>
    </React.StrictMode>
);

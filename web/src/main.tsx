import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';
import { QueryProvider } from '@/app/providers/QueryProvider';
import { ThemeProvider } from '@/app/providers/ThemeProvider';
import { router } from '@/app/router';
import '@/app/styles/globals.css';
import { AuthProvider } from '@/features/auth/AuthProvider';
import { ErrorBoundary } from '@/app/errors/ErrorBoundary';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
    <React.StrictMode>
      <ThemeProvider>
        <QueryProvider>
          <AuthProvider>
            <ErrorBoundary>
              <RouterProvider router={router} />
            </ErrorBoundary>
          </AuthProvider>
        </QueryProvider>
      </ThemeProvider>
    </React.StrictMode>
);

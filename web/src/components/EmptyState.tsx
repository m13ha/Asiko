import { ReactNode } from 'react';

export function EmptyState({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <div className={`grid justify-items-center text-center gap-2 p-6 border border-dashed border-gray-300 rounded-lg bg-gray-50 ${className}`}>
      {children}
    </div>
  );
}

export function EmptyTitle({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <h3 className={`text-lg font-semibold text-gray-900 mt-2 mb-0 ${className}`}>
      {children}
    </h3>
  );
}

export function EmptyDescription({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <p className={`text-sm text-gray-600 m-0 ${className}`}>
      {children}
    </p>
  );
}

export function EmptyAction({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <div className={`mt-2 ${className}`}>
      {children}
    </div>
  );
}
import { ReactNode } from 'react';

export interface CardProps {
  children: ReactNode;
  className?: string;
}

export function Card({ className = '', children }: CardProps) {
  return (
    <div className={`rounded-xl border border-[var(--border)] bg-[var(--bg-elevated)] text-[var(--text)] shadow-[var(--elev-1)] p-3 sm:p-6 ${className}`}>
      {children}
    </div>
  );
}

export function CardHeader({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <div className={`flex items-center justify-between mb-4 ${className}`}>
      {children}
    </div>
  );
}

export function CardTitle({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <h2 className={`text-lg font-semibold text-[var(--text)] m-0 ${className}`}>
      {children}
    </h2>
  );
}

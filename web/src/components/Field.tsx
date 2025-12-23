import { ReactNode } from 'react';

export function Field({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <div className={`w-full grid gap-1.5 ${className}`}>
      {children}
    </div>
  );
}

export function FieldLabel({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <label className={`text-xs text-gray-700 font-medium ${className}`}>
      {children}
    </label>
  );
}

export function FieldRow({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <div className={`w-full relative flex flex-wrap items-center gap-2 sm:gap-1.5 ${className}`}>
      {children}
    </div>
  );
}

export function IconSlot({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <span className={`absolute left-2.5 top-1/2 -translate-y-1/2 flex items-center justify-center w-4 h-4 text-gray-400 pointer-events-none ${className}`}>
      {children}
    </span>
  );
}

export function FieldError({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <span className={`text-xs text-red-500 font-medium ${className}`}>
      {children}
    </span>
  );
}
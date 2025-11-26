import { ReactNode } from 'react';

type BadgeTone = 'default' | 'success' | 'warning' | 'danger' | 'info' | 'muted';

export interface BadgeProps {
  children: ReactNode;
  tone?: BadgeTone;
  className?: string;
}

export function Badge({ children, tone = 'default', className = '' }: BadgeProps) {
  const toneClasses = {
    default: 'bg-gray-100 text-gray-800 border-gray-200',
    success: 'bg-green-100 text-green-800 border-green-200',
    warning: 'bg-yellow-100 text-yellow-800 border-yellow-200',
    danger: 'bg-red-100 text-red-800 border-red-200',
    info: 'bg-blue-100 text-blue-800 border-blue-200',
    muted: 'bg-gray-50 text-gray-600 border-gray-100'
  };

  return (
    <span className={`inline-flex items-center px-2 py-1 text-xs font-semibold rounded-full border ${toneClasses[tone]} ${className}`}>
      {children}
    </span>
  );
}
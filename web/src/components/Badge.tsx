import { ReactNode } from 'react';
import { Tag, TagProps } from 'primereact/tag';

type BadgeTone = 'default' | 'success' | 'warning' | 'danger' | 'info' | 'muted';

export interface BadgeProps extends Omit<TagProps, 'severity'> {
  children: ReactNode;
  tone?: BadgeTone;
}

export function Badge({ children, tone = 'default', className = '', ...props }: BadgeProps) {
  const getSeverity = () => {
    switch (tone) {
      case 'success': return 'success';
      case 'warning': return 'warning';
      case 'danger': return 'danger';
      case 'info': return 'info';
      case 'muted': return 'secondary';
      default: return undefined;
    }
  };

  return (
    <Tag
      severity={getSeverity()}
      className={className}
      {...props}
    >
      {children}
    </Tag>
  );
}
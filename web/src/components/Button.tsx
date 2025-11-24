import { forwardRef } from 'react';
import type { ButtonProps } from 'primereact/button';
import { Button as PrimeButton } from 'primereact/button';

type Variant = 'primary' | 'ghost';

export type AppButtonProps = Omit<ButtonProps, 'severity' | 'text' | 'link'> & {
  variant?: Variant;
};

export const Button = forwardRef<any, AppButtonProps>(function Button(
  { variant = 'primary', className, ...props },
  ref
) {
  const computedClass = [
    'app-button',
    variant === 'ghost' ? 'app-button--ghost' : 'app-button--primary',
    className,
  ]
    .filter(Boolean)
    .join(' ');
  return (
    <PrimeButton
      ref={ref}
      {...props}
      className={computedClass || undefined}
    />
  );
});

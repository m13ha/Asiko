import { forwardRef } from 'react';
import type { InputTextProps } from 'primereact/inputtext';
import { InputText } from 'primereact/inputtext';

export const Input = forwardRef<HTMLInputElement, InputTextProps>(function Input(
  { className, ...props },
  ref
) {
  const computedClass = ['app-input', className].filter(Boolean).join(' ');
  return <InputText ref={ref} {...props} className={computedClass || undefined} />;
});

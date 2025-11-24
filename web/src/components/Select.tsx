import { forwardRef } from 'react';
import type { DropdownProps } from 'primereact/dropdown';
import { Dropdown } from 'primereact/dropdown';

export type SelectProps = DropdownProps;

export const Select = forwardRef<any, SelectProps>(function Select(
  { className, ...props },
  ref
) {
  const computedClass = ['app-select', className].filter(Boolean).join(' ');
  return <Dropdown ref={ref} {...props} className={computedClass || undefined} />;
});

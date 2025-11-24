import { forwardRef } from 'react';
import type { InputTextareaProps } from 'primereact/inputtextarea';
import { InputTextarea } from 'primereact/inputtextarea';

export const Textarea = forwardRef<HTMLTextAreaElement, InputTextareaProps>(function Textarea(
  { className, autoResize = false, ...props },
  ref
) {
  const computedClass = ['app-textarea', className].filter(Boolean).join(' ');
  return (
    <InputTextarea
      ref={ref}
      {...props}
      className={computedClass || undefined}
      autoResize={autoResize}
    />
  );
});

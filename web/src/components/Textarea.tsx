import { forwardRef } from 'react';
import { InputTextarea, InputTextareaProps } from 'primereact/inputtextarea';

export const Textarea = forwardRef<HTMLTextAreaElement, InputTextareaProps>(function Textarea(
  { className, autoResize = false, ...props },
  ref
) {
  return (
    <InputTextarea
      ref={ref}
      {...props}
      className={className}
      autoResize={autoResize}
    />
  );
});

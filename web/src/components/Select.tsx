import { forwardRef } from 'react';
import { ChevronDown } from 'lucide-react';

export interface SelectOption {
  label: string;
  value: any;
}

export interface SelectProps extends Omit<React.SelectHTMLAttributes<HTMLSelectElement>, 'onChange'> {
  options: SelectOption[];
  optionLabel?: string;
  optionValue?: string;
  onChange?: (event: { value: any }) => void;
  error?: boolean;
  placeholder?: string;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(function Select(
  { options, optionLabel = 'label', optionValue = 'value', onChange, value, placeholder, className = '', error, ...props },
  ref
) {
  const baseClasses = 'w-full px-3 py-2.5 bg-[var(--bg-elevated)] border rounded-lg text-sm appearance-none transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 m-0 pr-10';
  const errorClasses = error 
    ? 'border-red-300 focus:border-red-500 focus:ring-red-500' 
    : 'border-gray-300 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600';

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const val = e.target.value;
    // If optionValue is value, val is always a string from the HTML element.
    // However, options might have numeric or enum values.
    // Find the original option to emit the correct type.
    const selectedOption = options.find(opt => String((opt as any)[optionValue]) === val);
    onChange?.({ value: selectedOption ? (selectedOption as any)[optionValue] : val });
  };

  return (
    <div className="relative w-full">
      <select
        ref={ref}
        value={value !== undefined ? String(value) : ''}
        onChange={handleChange}
        className={`${baseClasses} ${errorClasses} ${className}`}
        {...props}
      >
        {placeholder && <option value="" disabled>{placeholder}</option>}
        {options.map((option, idx) => (
          <option key={idx} value={String((option as any)[optionValue])}>
            {(option as any)[optionLabel]}
          </option>
        ))}
      </select>
      <div className="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--text-muted)] pointer-events-none">
        <ChevronDown size={16} />
      </div>
    </div>
  );
});

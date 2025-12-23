import { useEffect, useRef, useState } from 'react';
import { Button } from '@/components/Button';

export type StatusOption<T extends string | number> = {
  label: string;
  value: T;
};

interface StatusFilterProps<T extends string | number> {
  options: Array<StatusOption<T>>;
  selected: T[];
  onChange: (next: T[]) => void;
  label?: string;
  primaryCount?: number;
}

export function StatusFilter<T extends string | number>({
  options,
  selected,
  onChange,
  label = 'Status',
  primaryCount = 4,
}: StatusFilterProps<T>) {
  const isAllActive = selected.length === 0;
  const primary = options.slice(0, primaryCount);
  const overflow = options.slice(primaryCount);
  const overflowSelected = overflow.filter((opt) => selected.includes(opt.value)).length;
  const [open, setOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement | null>(null);

  const toggle = (value: T) => {
    if (selected.includes(value)) {
      onChange(selected.filter((item) => item !== value));
      return;
    }
    onChange([...selected, value]);
  };

  useEffect(() => {
    if (!open) return;
    const handleClick = (event: MouseEvent) => {
      if (!menuRef.current || !event.target) return;
      if (!menuRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, [open]);

  return (
    <div className="flex items-center gap-2 flex-wrap">
      <div className="text-xs uppercase tracking-wide text-[var(--text-muted)] mr-1">{label}</div>
      <div className="flex items-center gap-2 flex-wrap">
        <Button
          size="sm"
          variant={isAllActive ? 'primary' : 'outline'}
          onClick={() => onChange([])}
        >
          All
        </Button>
        {primary.map((option) => {
          const active = selected.includes(option.value);
          return (
            <Button
              key={String(option.value)}
              size="sm"
              variant={active ? 'primary' : 'outline'}
              onClick={() => toggle(option.value)}
            >
              {option.label}
            </Button>
          );
        })}
        {overflow.length > 0 && (
          <div className="relative" ref={menuRef}>
            <Button
              size="sm"
              variant={overflowSelected > 0 ? 'primary' : 'outline'}
              onClick={() => setOpen((prev) => !prev)}
            >
              More{overflowSelected > 0 ? ` (${overflowSelected})` : ''}
            </Button>
            {open && (
              <div className="absolute z-10 mt-2 w-44 rounded-lg border border-[var(--border)] bg-[var(--bg-elevated)] shadow-[var(--elev-2)] p-2">
                {overflow.map((option) => {
                  const active = selected.includes(option.value);
                  return (
                    <label key={String(option.value)} className="flex items-center gap-2 px-2 py-1 text-xs cursor-pointer">
                      <input
                        type="checkbox"
                        checked={active}
                        onChange={() => toggle(option.value)}
                      />
                      <span>{option.label}</span>
                    </label>
                  );
                })}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

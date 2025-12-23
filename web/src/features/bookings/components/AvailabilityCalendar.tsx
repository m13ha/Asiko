import { useMemo } from 'react';
import { Calendar } from 'primereact/calendar';
import { parseISO, format, isEqual } from 'date-fns';

type AvailabilityCalendarProps = {
  availableDates: string[];
  selectedDate?: string;
  onSelect: (date: string) => void;
};

type DateTemplateEvent = {
  day: number;
  month: number;
  year: number;
  otherMonth: boolean;
  today: boolean;
  selectable: boolean;
};

export function AvailabilityCalendar({ availableDates, selectedDate, onSelect }: AvailabilityCalendarProps) {
  const availableSet = useMemo(() => new Set(availableDates), [availableDates]);
  const selected = selectedDate ? parseISO(selectedDate) : null;
  const defaultDate = selected ?? (availableDates.length ? parseISO(availableDates[0]) : undefined);

  const handleChange = (event: { value?: Date | Date[] | null }) => {
    const value = (Array.isArray(event.value) ? event.value[0] : event.value) as Date | null;
    if (!value) return;
    const iso = format(value, 'yyyy-MM-dd');
    if (availableSet.has(iso)) {
      onSelect(iso);
    }
  };

  const dateTemplate = (event: DateTemplateEvent) => {
    const current = new Date(event.year, event.month, event.day);
    const iso = format(current, 'yyyy-MM-dd');
    const isAvailable = availableSet.has(iso);
    const isSelected = selected ? isEqual(current, selected) : false;
    const classes = [
      'inline-flex items-center justify-center w-9 h-9 rounded-md border text-sm font-medium transition',
      isAvailable ? 'text-[var(--text)] border-[var(--border)]' : 'text-[var(--text-muted)] border-transparent line-through opacity-50',
      isSelected ? 'bg-[var(--primary)] text-[var(--primary-contrast)] border-[var(--primary)] shadow-[0_0_0_2px_color-mix(in_oklab,var(--primary)_35%,transparent)]' : '',
      event.otherMonth ? 'opacity-40' : '',
    ];
    return (
      <span className={classes.join(' ')}>
        {event.day}
      </span>
    );
  };

  return (
    <div className="border border-[var(--border)] rounded-2xl p-3 bg-[var(--bg-elevated)]">
      <Calendar
        value={selected ?? defaultDate ?? null}
        onChange={handleChange}
        inline
        numberOfMonths={2}
        dateTemplate={dateTemplate}
        showWeek={false}
        className="w-full [&_.p-datepicker]:w-full [&_.p-datepicker]:!border-0 [&_.p-datepicker]:!shadow-none [&_.p-datepicker]:bg-[var(--bg-elevated)] [&_.p-datepicker]:text-[var(--text)] [&_.p-datepicker-header]:border-b [&_.p-datepicker-header]:border-[var(--border)]"
        prevIcon="pi pi-chevron-left"
        nextIcon="pi pi-chevron-right"
      />
    </div>
  );
}

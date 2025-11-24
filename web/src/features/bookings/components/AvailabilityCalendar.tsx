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
    const classes = ['availability-day'];
    if (isAvailable) classes.push('is-available');
    if (isSelected) classes.push('is-selected');
    if (event.otherMonth) classes.push('is-dimmed');
    if (!isAvailable) classes.push('is-disabled');
    return (
      <span className={classes.join(' ')}>
        {event.day}
      </span>
    );
  };

  return (
    <div className="availability-calendar">
      <Calendar
        value={selected ?? defaultDate ?? null}
        onChange={handleChange}
        inline
        numberOfMonths={2}
        dateTemplate={dateTemplate}
        showWeek={false}
        className="availability-calendar-widget"
        prevIcon="pi pi-chevron-left"
        nextIcon="pi pi-chevron-right"
      />
    </div>
  );
}

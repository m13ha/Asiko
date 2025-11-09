import { useEffect, useMemo, useState } from 'react';
import { addMonths, eachDayOfInterval, endOfMonth, endOfWeek, format, isSameDay, isSameMonth, parseISO, startOfMonth, startOfWeek } from 'date-fns';
import { Button } from '@/components/Button';
import { ChevronLeft, ChevronRight } from 'lucide-react';

type AvailabilityCalendarProps = {
  availableDates: string[];
  selectedDate?: string;
  onSelect: (date: string) => void;
};

const WEEKDAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

function toDate(date: string | undefined) {
  return date ? parseISO(date) : new Date();
}

export function AvailabilityCalendar({ availableDates, selectedDate, onSelect }: AvailabilityCalendarProps) {
  const selected = selectedDate ? toDate(selectedDate) : undefined;
  const initialFocus = selected ?? (availableDates.length ? toDate(availableDates[0]) : new Date());
  const [currentMonth, setCurrentMonth] = useState(startOfMonth(initialFocus));

  useEffect(() => {
    if (selected) {
      setCurrentMonth(startOfMonth(selected));
    }
  }, [selected]);

  useEffect(() => {
    if (!selected && availableDates.length) {
      setCurrentMonth(startOfMonth(toDate(availableDates[0])));
    }
  }, [availableDates, selected]);

  const availableSet = useMemo(() => new Set(availableDates), [availableDates]);
  const monthStart = startOfMonth(currentMonth);
  const monthEnd = endOfMonth(monthStart);
  const calendarStart = startOfWeek(monthStart, { weekStartsOn: 0 });
  const calendarEnd = endOfWeek(monthEnd, { weekStartsOn: 0 });
  const days = eachDayOfInterval({ start: calendarStart, end: calendarEnd });

  return (
    <div style={{ border: '1px solid var(--border)', borderRadius: 12, padding: 16, display: 'grid', gap: 12 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Button variant="ghost" onClick={() => setCurrentMonth((prev) => startOfMonth(addMonths(prev, -1)))}>
          <ChevronLeft size={16} /> Prev
        </Button>
        <strong>{format(monthStart, 'MMMM yyyy')}</strong>
        <Button variant="ghost" onClick={() => setCurrentMonth((prev) => startOfMonth(addMonths(prev, 1)))}>
          Next <ChevronRight size={16} />
        </Button>
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 1fr)', gap: 4, textAlign: 'center', fontSize: 12, textTransform: 'uppercase', color: 'var(--text-muted)' }}>
        {WEEKDAYS.map((day) => (<span key={day}>{day}</span>))}
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 1fr)', gap: 6 }}>
        {days.map((day) => {
          const key = format(day, 'yyyy-MM-dd');
          const isAvailable = availableSet.has(key);
          const isSelected = selected ? isSameDay(day, selected) : false;
          const disabled = !isAvailable;
          const isDimmed = !isSameMonth(day, monthStart);
          return (
            <button
              key={key}
              onClick={() => isAvailable && onSelect(key)}
              disabled={disabled}
              style={{
                height: 40,
                borderRadius: 10,
                border: '1px solid',
                borderColor: isSelected ? 'var(--primary)' : 'var(--border)',
                background: isSelected ? 'var(--primary)' : 'transparent',
                color: isSelected ? 'var(--primary-contrast)' : isDimmed ? 'var(--text-muted)' : 'var(--text)',
                opacity: disabled && !isSelected ? 0.4 : 1,
                cursor: disabled ? 'not-allowed' : 'pointer',
                fontWeight: isSelected ? 600 : 500,
              }}
            >
              {format(day, 'd')}
            </button>
          );
        })}
      </div>
    </div>
  );
}

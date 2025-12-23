import { useState } from 'react';
import { format, parseISO } from 'date-fns';
import { Button } from '@/components/Button';

interface GroupedDayPickerProps {
  availableDates: string[];
  selectedDate: string | null;
  onSelect: (date: string) => void;
}

export function GroupedDayPicker({ availableDates, selectedDate, onSelect }: GroupedDayPickerProps) {
  const [openMonth, setOpenMonth] = useState<string | null>(() => {
    if (selectedDate) {
      try {
        return format(parseISO(`${selectedDate}T00:00:00Z`), 'MMMM yyyy');
      } catch { return null; }
    }
    return null;
  });

  const groupedDates = availableDates.reduce((acc, dateStr) => {
    try {
      // Treat date as UTC to avoid timezone shifts
      const date = parseISO(`${dateStr}T00:00:00Z`);
      const month = format(date, 'MMMM yyyy');
      if (!acc[month]) {
        acc[month] = [];
      }
      acc[month].push(dateStr);
      return acc;
    } catch (e) {
      console.error(`Could not parse date: ${dateStr}`, e);
      return acc;
    }
  }, {} as Record<string, string[]>);

  const toggleMonth = (month: string) => {
    setOpenMonth(openMonth === month ? null : month);
  };

  // Automatically open the first month if none is selected
  if (openMonth === null && Object.keys(groupedDates).length > 0) {
    setOpenMonth(Object.keys(groupedDates)[0]);
  }

  return (
    <div className="grid gap-2">
      {Object.entries(groupedDates).map(([month, dates]) => (
        <div key={month} className="rounded-xl border border-[var(--border)] bg-[var(--bg-elevated)]">
          <button
            onClick={() => toggleMonth(month)}
            className="w-full text-left font-semibold px-4 py-3 flex justify-between items-center"
            aria-expanded={openMonth === month}
          >
            <span>{month}</span>
            <i className={`pi pi-chevron-down transition-transform ${openMonth === month ? 'rotate-180' : ''}`} aria-hidden="true" />
          </button>
          {openMonth === month && (
            <div className="grid grid-cols-3 sm:grid-cols-4 gap-2 p-3 border-t border-[var(--border)]">
              {dates.map((dateStr) => {
                const date = parseISO(`${dateStr}T00:00:00Z`);
                const isSelected = dateStr === selectedDate;
                return (
                  <Button
                    key={dateStr}
                    variant={isSelected ? 'primary' : 'outline'}
                    onClick={() => onSelect(dateStr)}
                    className="flex flex-col items-center justify-center h-20 text-center"
                  >
                    <span className="text-xs font-medium">{format(date, 'EEE')}</span>
                    <span className="text-2xl font-bold">{format(date, 'd')}</span>
                  </Button>
                );
              })}
            </div>
          )}
        </div>
      ))}
    </div>
  );
}

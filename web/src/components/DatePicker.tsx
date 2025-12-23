import { forwardRef, useState, useRef, useEffect, useMemo } from 'react';
import { createPortal } from 'react-dom';
import { Calendar as CalendarIcon, X, ChevronLeft, ChevronRight } from 'lucide-react';
import { 
  format, 
  addMonths, 
  subMonths, 
  startOfMonth, 
  endOfMonth, 
  startOfWeek, 
  endOfWeek, 
  isSameMonth, 
  isSameDay, 
  addDays, 
  eachDayOfInterval,
  isBefore,
  startOfDay
} from 'date-fns';
import { Input } from './Input';
import { Button } from './Button';

type DatePickerProps = {
  value?: Date | string | null;
  onChange?: (value: Date | null) => void;
  placeholder?: string;
  disabled?: boolean;
  minDate?: Date | string;
  maxDate?: Date | string;
  className?: string;
};

export const DatePicker = forwardRef<HTMLInputElement, DatePickerProps>(function DatePicker(
  { value, onChange, placeholder = 'Select date', disabled, minDate, maxDate, className = '' },
  ref
) {
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const triggerRef = useRef<HTMLDivElement>(null);

  const parsedValue = useMemo(() => {
    if (!value) return null;
    const d = new Date(value);
    return isNaN(d.getTime()) ? null : d;
  }, [value]);

  const [viewDate, setViewDate] = useState(parsedValue || new Date());

  useEffect(() => {
    if (parsedValue) {
      setViewDate(parsedValue);
    }
  }, [parsedValue]);

  const openPicker = () => {
    if (!disabled) {
      if (parsedValue) setViewDate(parsedValue);
      setIsOpen(true);
    }
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node) &&
        triggerRef.current &&
        !triggerRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen]);

  const handleDateSelect = (date: Date) => {
    onChange?.(date);
    setIsOpen(false);
  };

  const renderHeader = () => (
    <div className="flex items-center justify-between mb-4">
      <button 
        type="button"
        onClick={() => setViewDate(subMonths(viewDate, 1))}
        className="p-2 hover:bg-[var(--bg-muted)] rounded-full transition-colors text-[var(--primary)]"
      >
        <ChevronLeft size={20} />
      </button>
      <span className="font-bold text-[var(--text)]">
        {format(viewDate, 'MMMM yyyy')}
      </span>
      <button 
        type="button"
        onClick={() => setViewDate(addMonths(viewDate, 1))}
        className="p-2 hover:bg-[var(--bg-muted)] rounded-full transition-colors text-[var(--primary)]"
      >
        <ChevronRight size={20} />
      </button>
    </div>
  );

  const renderDays = () => {
    const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
    return (
      <div className="grid grid-cols-7 mb-2">
        {days.map(day => (
          <div key={day} className="text-center text-[10px] font-bold text-[var(--text-muted)] uppercase tracking-wider">
            {day}
          </div>
        ))}
      </div>
    );
  };

  const renderCells = () => {
    const monthStart = startOfMonth(viewDate);
    const monthEnd = endOfMonth(monthStart);
    const startDate = startOfWeek(monthStart);
    const endDate = endOfWeek(monthEnd);

    const dateFormat = "d";
    const rows = [];

    const days = eachDayOfInterval({ start: startDate, end: endDate });

    const min = minDate ? startOfDay(new Date(minDate)) : null;
    const max = maxDate ? startOfDay(new Date(maxDate)) : null;

    return (
      <div className="grid grid-cols-7 gap-1">
        {days.map((day, idx) => {
          const isSelected = parsedValue && isSameDay(day, parsedValue);
          const isCurrentMonth = isSameMonth(day, monthStart);
          const isDisabled = (min && isBefore(day, min)) || (max && isBefore(max, day));

          return (
            <button
              key={idx}
              type="button"
              disabled={isDisabled}
              onClick={() => handleDateSelect(day)}
              className={`
                aspect-square flex items-center justify-center text-sm rounded-lg transition-all
                ${!isCurrentMonth ? 'text-[var(--text-muted)] opacity-30' : 'text-[var(--text)]'}
                ${isSelected ? 'bg-[var(--primary)] text-white font-bold scale-105 shadow-md' : 'hover:bg-[var(--primary)]/10'}
                ${isDisabled ? 'opacity-20 cursor-not-allowed grayscale' : 'cursor-pointer'}
              `}
            >
              {format(day, dateFormat)}
            </button>
          );
        })}
      </div>
    );
  };

  const pickerContent = (
    <div
      ref={containerRef}
      className={`
        fixed inset-x-0 bottom-0 z-[100] bg-[var(--bg-elevated)] border-t border-[var(--border)] rounded-t-3xl shadow-2xl p-6 transition-transform duration-300 ease-in-out sm:absolute sm:inset-auto sm:mt-2 sm:rounded-xl sm:border sm:w-80 sm:shadow-lg
        ${isOpen ? 'translate-y-0 opacity-100' : 'translate-y-full opacity-0 pointer-events-none sm:translate-y-0 sm:opacity-0'}
      `}
      style={window.innerWidth >= 640 && triggerRef.current ? {
        top: triggerRef.current.getBoundingClientRect().bottom + window.scrollY,
        left: triggerRef.current.getBoundingClientRect().left + window.scrollX
      } : {}}
    >
      <div className="flex items-center justify-between mb-4 sm:hidden">
        <h3 className="text-lg font-bold text-[var(--text)]">Select Date</h3>
        <button onClick={() => setIsOpen(false)} className="p-2 text-[var(--text-muted)]">
          <X size={24} />
        </button>
      </div>

      {renderHeader()}
      {renderDays()}
      {renderCells()}

      <div className="mt-6 sm:hidden">
        <Button variant="outline" className="w-full" onClick={() => setIsOpen(false)}>Cancel</Button>
      </div>
    </div>
  );

  return (
    <div className={`relative ${className}`} ref={triggerRef}>
      <Input
        ref={ref}
        value={parsedValue ? format(parsedValue, 'PPP') : ''}
        readOnly
        placeholder={placeholder}
        disabled={disabled}
        onClick={openPicker}
        className="cursor-pointer pl-10"
        onFocus={(e) => {
          e.currentTarget.blur();
          openPicker();
        }}
      />
      <div className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--text-muted)] pointer-events-none">
        <CalendarIcon size={16} />
      </div>
      
      {isOpen && createPortal(
        <>
          <div 
            className="fixed inset-0 bg-black/40 z-[90] sm:hidden transition-opacity" 
            onClick={() => setIsOpen(false)}
          />
          {pickerContent}
        </>,
        document.body
      )}
    </div>
  );
});

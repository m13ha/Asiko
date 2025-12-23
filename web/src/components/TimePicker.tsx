import { forwardRef, useState, useRef, useEffect, useMemo } from 'react';
import { createPortal } from 'react-dom';
import { Clock, X, ChevronUp, ChevronDown } from 'lucide-react';
import { format, setHours, setMinutes, parse } from 'date-fns';
import { Input } from './Input';
import { Button } from './Button';

type TimePickerProps = {
  value?: Date | null;
  onChange?: (value: Date | null) => void;
  placeholder?: string;
  disabled?: boolean;
  minTime?: Date;
  className?: string;
};

export const TimePicker = forwardRef<HTMLInputElement, TimePickerProps>(function TimePicker(
  { value, onChange, placeholder = 'Select time', disabled, minTime, className = '' },
  ref
) {
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const triggerRef = useRef<HTMLDivElement>(null);

  // Default to 12:00 PM if no value
  const displayDate = useMemo(() => {
    if (value instanceof Date) return value;
    const d = new Date();
    d.setHours(12, 0, 0, 0);
    return d;
  }, [value]);
  
  const [selectedHour, setSelectedHour] = useState(parseInt(format(displayDate, 'h')));
  const [selectedMinute, setSelectedMinute] = useState(Math.round(displayDate.getMinutes() / 5) * 5);
  const [period, setPeriod] = useState(format(displayDate, 'a').toUpperCase() as 'AM' | 'PM');

  const syncLocalState = (date: Date) => {
    setSelectedHour(parseInt(format(date, 'h')));
    setSelectedMinute(Math.round(date.getMinutes() / 5) * 5);
    setPeriod(format(date, 'a').toUpperCase() as 'AM' | 'PM');
  };

  useEffect(() => {
    if (value) {
      syncLocalState(value);
    }
  }, [value]);

  const openPicker = () => {
    if (!disabled) {
      if (value) syncLocalState(value);
      setIsOpen(true);
    }
  };

  const handleApply = () => {
    let hour = selectedHour;
    if (period === 'PM' && hour !== 12) hour += 12;
    if (period === 'AM' && hour === 12) hour = 0;

    const newTime = setMinutes(setHours(new Date(), hour), selectedMinute);
    newTime.setSeconds(0, 0);
    
    onChange?.(newTime);
    setIsOpen(false);
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

  const hours = Array.from({ length: 12 }, (_, i) => i + 1);
  const minutes = Array.from({ length: 12 }, (_, i) => i * 5);

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
      <div className="flex items-center justify-between mb-6 sm:hidden">
        <h3 className="text-lg font-bold text-[var(--text)]">Select Time</h3>
        <button onClick={() => setIsOpen(false)} className="p-2 text-[var(--text-muted)]">
          <X size={24} />
        </button>
      </div>

      <div className="flex justify-center items-center gap-4 mb-8">
        {/* Hours */}
        <div className="flex flex-col items-center">
          <button onClick={() => setSelectedHour(h => h === 12 ? 1 : h + 1)} className="p-2 text-[var(--primary)]"><ChevronUp /></button>
          <span className="text-4xl font-mono font-bold w-16 text-center">{selectedHour.toString().padStart(2, '0')}</span>
          <button onClick={() => setSelectedHour(h => h === 1 ? 12 : h - 1)} className="p-2 text-[var(--primary)]"><ChevronDown /></button>
        </div>
        
        <span className="text-4xl font-bold pb-2">:</span>

        {/* Minutes */}
        <div className="flex flex-col items-center">
          <button onClick={() => setSelectedMinute(m => (m + 5) % 60)} className="p-2 text-[var(--primary)]"><ChevronUp /></button>
          <span className="text-4xl font-mono font-bold w-16 text-center">{selectedMinute.toString().padStart(2, '0')}</span>
          <button onClick={() => setSelectedMinute(m => (m - 5 + 60) % 60)} className="p-2 text-[var(--primary)]"><ChevronDown /></button>
        </div>

        {/* AM/PM */}
        <div className="flex flex-col gap-2 ml-2">
          <button 
            type="button"
            onClick={() => setPeriod('AM')}
            className={`px-3 py-1.5 rounded-md font-bold text-xs ${period === 'AM' ? 'bg-[var(--primary)] text-white' : 'bg-[var(--bg-muted)] text-[var(--text-muted)]'}`}
          >AM</button>
          <button 
            type="button"
            onClick={() => setPeriod('PM')}
            className={`px-3 py-1.5 rounded-md font-bold text-xs ${period === 'PM' ? 'bg-[var(--primary)] text-white' : 'bg-[var(--bg-muted)] text-[var(--text-muted)]'}`}
          >PM</button>
        </div>
      </div>

      <div className="flex gap-3">
        <Button type="button" variant="outline" className="flex-1" onClick={() => setIsOpen(false)}>Cancel</Button>
        <Button type="button" variant="primary" className="flex-1" onClick={handleApply}>Apply</Button>
      </div>
    </div>
  );

  return (
    <div className={`relative ${className}`} ref={triggerRef}>
      <Input
        ref={ref}
        value={value ? format(value, 'hh:mm a') : ''}
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
        <Clock size={16} />
      </div>
      
      {isOpen && createPortal(
        <>
          {/* Mobile Overlay */}
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

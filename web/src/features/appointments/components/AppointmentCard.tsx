import { useState } from 'react';
import { Calendar, Clock, Users, ChevronRight, Copy, Check } from 'lucide-react';
import { Link } from 'react-router-dom';
import { format } from 'date-fns';

function formatLabel(value?: string) {
  if (!value) return '';
  return value.charAt(0).toUpperCase() + value.slice(1);
}

function toDate(value?: string) {
  if (!value) return null;
  const d = new Date(value);
  return Number.isNaN(d.getTime()) ? null : d;
}

function formatDateRange(startDate?: string, endDate?: string) {
  const start = toDate(startDate);
  const end = toDate(endDate);
  if (start && end) {
    const sameDay = start.toDateString() === end.toDateString();
    return sameDay
      ? format(start, 'EEE, MMM d, yyyy')
      : `${format(start, 'MMM d, yyyy')} → ${format(end, 'MMM d, yyyy')}`;
  }
  if (start) return format(start, 'EEE, MMM d, yyyy');
  if (end) return format(end, 'EEE, MMM d, yyyy');
  return 'Date TBD';
}

function formatTimeRange(startTime?: string, endTime?: string) {
  const start = toDate(startTime);
  const end = toDate(endTime);
  if (start && end) return `${format(start, 'p')} – ${format(end, 'p')}`;
  if (start) return format(start, 'p');
  if (end) return format(end, 'p');
  return 'Time TBD';
}

function getStatusClasses(status?: string) {
  const value = (status || '').toLowerCase();
  if (['active', 'ongoing'].includes(value)) return { 
    bg: 'bg-green-100', 
    text: 'text-green-700', 
    border: 'border-green-500' 
  };
  if (['pending', 'draft'].includes(value)) return { 
    bg: 'bg-yellow-100', 
    text: 'text-yellow-700', 
    border: 'border-yellow-500' 
  };
  if (['completed'].includes(value)) return { 
    bg: 'bg-blue-100', 
    text: 'text-blue-700', 
    border: 'border-blue-500' 
  };
  if (['canceled', 'cancelled', 'expired'].includes(value)) return { 
    bg: 'bg-gray-100', 
    text: 'text-gray-600', 
    border: 'border-gray-400' 
  };
  return { 
    bg: 'bg-gray-100', 
    text: 'text-gray-600', 
    border: 'border-gray-400' 
  };
}

export function AppointmentCard({ item }: { item: any }) {
  const [copiedCode, setCopiedCode] = useState(false);
  const dateLabel = formatDateRange(item.startDate, item.endDate);
  const timeLabel = formatTimeRange(item.startTime, item.endTime);
  const statusClasses = getStatusClasses(item.status);

  const copyToClipboard = (code: string) => {
    navigator.clipboard.writeText(code);
    setCopiedCode(true);
    setTimeout(() => setCopiedCode(false), 2000);
  };

  return (
    <div className={`bg-white rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden border-l-4 ${statusClasses.border}`}>
      <div className="p-6">
        {/* Code Section - Prominent */}
        <div className="mb-4 p-3 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg border border-blue-200">
          <div className="flex items-center justify-between">
            <div className="flex-1 min-w-0">
              <p className="text-xs font-semibold text-gray-600 mb-1">CODE</p>
              <p className="text-lg font-bold font-mono text-gray-900 tracking-wide truncate">
                {item.appCode || '—'}
              </p>
            </div>
            {item.appCode && (
              <button
                onClick={() => copyToClipboard(item.appCode)}
                className={`ml-2 px-2 py-2 ${copiedCode ? 'bg-green-600' : 'bg-blue-600 hover:bg-blue-700'} text-white rounded-md text-xs font-semibold transition-all flex items-center gap-1 shadow-md`}
              >
                {copiedCode ? (
                  <>
                    <Check className="w-4 h-4" />
                    <span className="hidden sm:inline">Copied!</span>
                  </>
                ) : (
                  <>
                    <Copy className="w-4 h-4" />
                    <span className="hidden sm:inline">Copy</span>
                  </>
                )}
              </button>
            )}
          </div>
        </div>

        {/* Header Section */}
        <div className="mb-3">
          <div className="flex items-center gap-2 mb-2 flex-wrap">
            <span className={`px-2 py-1 rounded-full text-xs font-semibold ${statusClasses.bg} ${statusClasses.text}`}>
              {formatLabel(item.status)}
            </span>
            <span className="px-2 py-1 rounded-full text-xs font-semibold bg-blue-50 text-blue-700">
              {formatLabel(String(item.type))}
            </span>
          </div>
          <h2 className="text-lg font-bold text-gray-800 mb-1 line-clamp-2">
            {item.title || 'Untitled appointment'}
          </h2>
          {item.description && <p className="text-sm text-gray-600 line-clamp-2">{item.description}</p>}
        </div>

        {/* Details Grid */}
        <div className="space-y-2 mb-3">
          <div className="flex items-center gap-2 text-gray-700">
            <div className="p-1.5 bg-blue-50 rounded">
              <Calendar className="w-4 h-4 text-blue-600" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-xs text-gray-500 font-medium">Date</p>
              <p className="text-xs font-semibold truncate">{dateLabel}</p>
            </div>
          </div>

          <div className="flex items-center gap-2 text-gray-700">
            <div className="p-1.5 bg-purple-50 rounded">
              <Clock className="w-4 h-4 text-purple-600" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-xs text-gray-500 font-medium">Time</p>
              <p className="text-xs font-semibold truncate">{timeLabel}</p>
            </div>
          </div>

          {item.maxAttendees && (
            <div className="flex items-center gap-2 text-gray-700">
              <div className="p-1.5 bg-orange-50 rounded">
                <Users className="w-4 h-4 text-orange-600" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-xs text-gray-500 font-medium">Capacity</p>
                <p className="text-xs font-semibold">{item.maxAttendees} slots</p>
              </div>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between pt-3 border-t border-gray-100">
          <Link 
            to={`/appointments/${item.id}`} 
            state={{ appointment: item }}
            className="text-blue-600 hover:text-blue-700 font-semibold text-xs flex items-center gap-1 transition-colors"
          >
            <button className="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-md font-semibold text-xs transition-colors">
              Manage
            </button>
          </Link>
        </div>
      </div>
    </div>
  );
}
import { useState } from 'react';
import { Calendar, Clock, ChevronRight, Copy, Check, Mail, User } from 'lucide-react';
import { format } from 'date-fns';

interface BookingCardProps {
  booking: {
    id?: string;
    bookingCode?: string;
    appCode?: string;
    date?: string;
    startTime?: string;
    endTime?: string;
    status?: string;
    email?: string;
    name?: string;
    seatsBooked?: number;
  };
  onAction?: (action: string, booking: any) => void;
  showActions?: boolean;
}

function formatTime(timeStr?: string) {
  if (!timeStr) return 'Time TBD';
  try {
    const date = new Date(timeStr);
    return format(date, 'p');
  } catch {
    return timeStr;
  }
}

function formatDate(dateStr?: string) {
  if (!dateStr) return 'Date TBD';
  try {
    const date = new Date(dateStr);
    return format(date, 'EEE, MMM d, yyyy');
  } catch {
    return dateStr;
  }
}

function getStatusClasses(status?: string) {
  const s = (status || '').toLowerCase();
  if (['active', 'confirmed'].includes(s)) return { 
    bg: 'bg-green-100', 
    text: 'text-green-700', 
    border: 'border-green-500' 
  };
  if (['pending'].includes(s)) return { 
    bg: 'bg-yellow-100', 
    text: 'text-yellow-700', 
    border: 'border-yellow-500' 
  };
  if (['cancelled', 'canceled', 'rejected'].includes(s)) return { 
    bg: 'bg-red-100', 
    text: 'text-red-700', 
    border: 'border-red-500' 
  };
  return { 
    bg: 'bg-gray-100', 
    text: 'text-gray-600', 
    border: 'border-gray-400' 
  };
}

export function BookingCard({ booking, onAction, showActions = false }: BookingCardProps) {
  const [copiedCode, setCopiedCode] = useState<string | null>(null);
  const statusClasses = getStatusClasses(booking.status);

  const copyToClipboard = (code: string) => {
    navigator.clipboard.writeText(code);
    setCopiedCode(code);
    setTimeout(() => setCopiedCode(null), 2000);
  };

  return (
    <div className={`bg-white rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden border-l-4 ${statusClasses.border}`}>
      <div className="p-6">
        {/* Booking Code Section */}
        {booking.bookingCode && (
          <div className="mb-4 p-3 bg-gradient-to-r from-purple-50 to-pink-50 rounded-lg border border-purple-200">
            <div className="flex items-center justify-between">
              <div className="flex-1 min-w-0">
                <p className="text-xs font-semibold text-gray-600 mb-1">CODE</p>
                <p className="text-lg font-bold font-mono text-gray-900 tracking-wide truncate">
                  {booking.bookingCode}
                </p>
              </div>
              <button
                onClick={() => copyToClipboard(booking.bookingCode!)}
                className={`ml-2 px-2 py-2 ${copiedCode === booking.bookingCode ? 'bg-green-600' : 'bg-purple-600 hover:bg-purple-700'} text-white rounded-md text-xs font-semibold transition-all flex items-center gap-1 shadow-md`}
              >
                {copiedCode === booking.bookingCode ? (
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
            </div>
          </div>
        )}

        {/* Header Section */}
        <div className="mb-3">
          <div className="flex items-center gap-2 mb-2 flex-wrap">
            <span className={`px-2 py-1 rounded-full text-xs font-semibold ${statusClasses.bg} ${statusClasses.text}`}>
              {booking.status}
            </span>
            {booking.seatsBooked && (
              <span className="px-2 py-1 rounded-full text-xs font-semibold bg-blue-50 text-blue-700">
                {booking.seatsBooked} seat{booking.seatsBooked !== 1 ? 's' : ''}
              </span>
            )}
          </div>
          <h2 className="text-lg font-bold text-gray-800 mb-1 truncate">
            Booking for {booking.appCode}
          </h2>
        </div>

        {/* Details Grid */}
        <div className="space-y-2 mb-3">

          {booking.name && (
            <div className="flex items-center gap-2 text-gray-700">
              <div className="p-1.5 bg-green-50 rounded">
                <User className="w-4 h-4 text-green-600" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-xs text-gray-500 font-medium">Name</p>
                <p className="text-xs font-semibold truncate">{booking.name}</p>
              </div>
            </div>
          )}

          {booking.email && (
            <div className="flex items-center gap-2 text-gray-700">
              <div className="p-1.5 bg-orange-50 rounded">
                <Mail className="w-4 h-4 text-orange-600" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-xs text-gray-500 font-medium">Email</p>
                <p className="text-xs font-semibold truncate">{booking.email}</p>
              </div>
            </div>
          )}
          
          <div className="flex items-center gap-2 text-gray-700">
            <div className="p-1.5 bg-blue-50 rounded">
              <Calendar className="w-4 h-4 text-blue-600" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-xs text-gray-500 font-medium">Date</p>
              <p className="text-xs font-semibold truncate">{formatDate(booking.date)}</p>
            </div>
          </div>

          <div className="flex items-center gap-2 text-gray-700">
            <div className="p-1.5 bg-purple-50 rounded">
              <Clock className="w-4 h-4 text-purple-600" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-xs text-gray-500 font-medium">Time</p>
              <p className="text-xs font-semibold truncate">{formatTime(booking.startTime)} - {formatTime(booking.endTime)}</p>
            </div>
          </div>
        </div>

        {/* Appointment Code Section */}
        {/* {booking.appCode && (
          <div className="mb-3 p-2 bg-gray-50 rounded border">
            <div className="flex items-center justify-between">
              <div className="min-w-0 flex-1">
                <p className="text-xs text-gray-500 font-medium">Appointment</p>
                <p className="font-mono text-xs font-semibold text-gray-900 truncate">{booking.appCode}</p>
              </div>
              <button
                onClick={() => copyToClipboard(booking.appCode!)}
                className="p-1 text-gray-600 hover:text-gray-800 transition-colors ml-2"
              >
                {copiedCode === booking.appCode ? (
                  <Check className="w-3 h-3 text-green-600" />
                ) : (
                  <Copy className="w-3 h-3" />
                )}
              </button>
            </div>
          </div>
        )} */}

        {/* Footer */}
        {showActions && onAction && (
          <div className="flex items-center justify-between pt-3 border-t border-gray-100">
            <button
              onClick={() => onAction('view', booking)}
              className="text-blue-600 hover:text-blue-700 font-semibold text-xs flex items-center gap-1 transition-colors"
            >
              <span className="hidden sm:inline">View Details</span>
              <span className="sm:hidden">View</span>
              <ChevronRight className="w-3 h-3" />
            </button>
            <div className="flex gap-1">
              {booking.status === 'active' && (
                <>
                  <button
                    onClick={() => onAction('update', booking)}
                    className="px-2 py-1 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded text-xs font-semibold transition-colors"
                  >
                    Update
                  </button>
                  <button
                    onClick={() => onAction('cancel', booking)}
                    className="px-2 py-1 bg-red-600 hover:bg-red-700 text-white rounded text-xs font-semibold transition-colors"
                  >
                    Cancel
                  </button>
                </>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
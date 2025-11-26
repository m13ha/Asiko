import { useMyBookings } from '../hooks';
import { BookingCard } from '@/components/BookingCard';
import { PaginatedGrid } from '@/components/PaginatedGrid';
import { usePagination } from '@/hooks/usePagination';
import { EmptyState, EmptyTitle, EmptyDescription } from '@/components/EmptyState';

export function MyBookingsPage() {
  const pagination = usePagination(1, 10);
  const { data, isLoading, error } = useMyBookings(pagination.params);
  
  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8 bg-white rounded-2xl shadow-lg p-6">
          <h1 className="text-3xl font-bold text-gray-800">My Bookings</h1>
          <p className="text-gray-600 mt-2">Manage your appointment bookings</p>
        </div>
        
        <PaginatedGrid
          data={data}
          isLoading={isLoading}
          error={error}
          onPageChange={pagination.updatePage}
          renderItem={(booking: any) => (
            <BookingCard 
              key={booking.id} 
              booking={booking}
              showActions={true}
              onAction={(action, booking) => {
                // Handle booking actions (view, update, cancel)
                console.log(action, booking);
              }}
            />
          )}
          emptyState={
            <div className="bg-white rounded-2xl shadow-lg p-12 text-center">
              <EmptyState>
                <EmptyTitle>No bookings yet</EmptyTitle>
                <EmptyDescription>When you book a slot, it will appear here.</EmptyDescription>
              </EmptyState>
            </div>
          }
        />
      </div>
    </div>
  );
}
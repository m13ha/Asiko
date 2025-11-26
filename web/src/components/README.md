# Reusable Pagination Components

## Overview
This directory contains reusable components for displaying paginated data from the Asiko API.

## Components

### PaginatedList
A generic component for displaying paginated lists with built-in pagination controls.

**Props:**
- `data` - Paginated response from API
- `isLoading` - Loading state
- `error` - Error state  
- `onPageChange` - Callback for page changes
- `renderItem` - Function to render each item
- `emptyState` - Component to show when no items
- `className` - Container CSS classes
- `itemsClassName` - Items grid CSS classes

**Usage:**
```tsx
<PaginatedList
  data={paginatedData}
  isLoading={isLoading}
  error={error}
  onPageChange={handlePageChange}
  renderItem={(item) => <ItemCard key={item.id} item={item} />}
  emptyState={<EmptyState>No items found</EmptyState>}
/>
```

### BookingCard
A reusable card component for displaying booking information.

**Props:**
- `booking` - Booking object with id, codes, dates, status, etc.
- `onAction` - Callback for action buttons (view, update, cancel)
- `showActions` - Whether to show action buttons

**Usage:**
```tsx
<BookingCard 
  booking={bookingData}
  showActions={true}
  onAction={(action, booking) => handleAction(action, booking)}
/>
```

## Hooks

### usePagination
Custom hook for managing pagination state.

**Returns:**
- `page` - Current page number
- `size` - Page size
- `updatePage` - Function to change page
- `updateSize` - Function to change page size
- `reset` - Function to reset to initial values
- `params` - Object with page/size for API calls

**Usage:**
```tsx
const pagination = usePagination(1, 10);
const { data } = useMyBookings(pagination.params);

// In component
<PaginatedList onPageChange={pagination.updatePage} />
```

## Updated API Integration

All paginated APIs now support `page` and `size` parameters:

- `useMyAppointments({ statuses, page, size })`
- `useMyBookings({ page, size })`
- `useNotifications({ page, size })`
- `useAvailableSlots(appCode, { page, size })`
- `useAvailableSlotsByDay(appCode, date, { page, size })`

## Migration

Existing pages have been updated to use these components:
- `MyAppointmentsPage` - Uses PaginatedList with AppointmentCard
- `MyBookingsPage` - Uses PaginatedList with BookingCard  
- `NotificationsPage` - Uses PaginatedList with notification cards

The components handle loading states, errors, empty states, and pagination automatically.
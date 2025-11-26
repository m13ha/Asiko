import { useState, useEffect } from 'react';
import { Button } from './Button';
import { Spinner } from './Spinner';

export interface PaginationInfo {
  page?: number;
  perPage?: number;
  total?: number;
  totalPages?: number;
}

export interface PaginatedData<T> extends PaginationInfo {
  items?: T[];
}

interface PaginatedListProps<T> {
  data?: PaginatedData<T>;
  isLoading?: boolean;
  error?: unknown;
  onPageChange?: (page: number) => void;
  renderItem: (item: T, index: number) => React.ReactNode;
  emptyState?: React.ReactNode;
  className?: string;
  itemsClassName?: string;
}

export function PaginatedList<T>({
  data,
  isLoading,
  error,
  onPageChange,
  renderItem,
  emptyState,
  className = 'min-h-screen bg-gradient-to-br from-blue-50 to-indigo-50',
  itemsClassName = 'space-y-6'
}: PaginatedListProps<T>) {
  const [currentPage, setCurrentPage] = useState(1);

  useEffect(() => {
    if (data?.page && data.page !== currentPage) {
      setCurrentPage(data.page);
    }
  }, [data?.page]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    onPageChange?.(page);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-8">
        <Spinner />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-red-600 py-4">
        Failed to load data. Please try again.
      </div>
    );
  }

  if (!data?.items?.length) {
    return emptyState || (
      <div className="text-gray-500 py-8 text-center">
        No items found.
      </div>
    );
  }

  return (
    <div className={className}>
      <div className="max-w-6xl mx-auto">
        <div className={itemsClassName}>
          {data.items.map((item, index) => renderItem(item, index))}
        </div>
        
        {data.totalPages && data.totalPages > 1 && (
          <Pagination
            currentPage={currentPage}
            totalPages={data.totalPages}
            onPageChange={handlePageChange}
          />
        )}
      </div>
    </div>
  );
}

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

function Pagination({ currentPage, totalPages, onPageChange }: PaginationProps) {
  const getVisiblePages = () => {
    const delta = 2;
    const range = [];
    const rangeWithDots = [];

    for (let i = Math.max(2, currentPage - delta); i <= Math.min(totalPages - 1, currentPage + delta); i++) {
      range.push(i);
    }

    if (currentPage - delta > 2) {
      rangeWithDots.push(1, '...');
    } else {
      rangeWithDots.push(1);
    }

    rangeWithDots.push(...range);

    if (currentPage + delta < totalPages - 1) {
      rangeWithDots.push('...', totalPages);
    } else if (totalPages > 1) {
      rangeWithDots.push(totalPages);
    }

    return rangeWithDots;
  };

  return (
    <div className="flex justify-center items-center gap-2 mt-6">
      <Button
        variant="outline"
        size="sm"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage <= 1}
      >
        Previous
      </Button>
      
      {getVisiblePages().map((page, index) => (
        <span key={index}>
          {page === '...' ? (
            <span className="px-2 text-gray-400">...</span>
          ) : (
            <Button
              variant={page === currentPage ? 'primary' : 'outline'}
              size="sm"
              onClick={() => onPageChange(page as number)}
            >
              {page}
            </Button>
          )}
        </span>
      ))}
      
      <Button
        variant="outline"
        size="sm"
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage >= totalPages}
      >
        Next
      </Button>
    </div>
  );
}
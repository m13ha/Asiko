import { useState, useEffect } from 'react';
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

interface PaginatedGridProps<T> {
  data?: PaginatedData<T>;
  isLoading?: boolean;
  error?: unknown;
  onPageChange?: (page: number) => void;
  renderItem: (item: T, index: number) => React.ReactNode;
  emptyState?: React.ReactNode;
}

export function PaginatedGrid<T>({
  data,
  isLoading,
  error,
  onPageChange,
  renderItem,
  emptyState
}: PaginatedGridProps<T>) {
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
    <div>
      <div className="grid gap-6 grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 mb-8">
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
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage <= 1}
        className="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-700 hover:bg-gray-50 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed transition-colors"
      >
        Previous
      </button>
      
      {getVisiblePages().map((page, index) => (
        <span key={index}>
          {page === '...' ? (
            <span className="px-2 text-gray-400">...</span>
          ) : (
            <button
              onClick={() => onPageChange(page as number)}
              className={`px-3 py-2 border rounded-lg transition-colors ${
                page === currentPage
                  ? 'bg-blue-600 text-white border-blue-600'
                  : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
              }`}
            >
              {page}
            </button>
          )}
        </span>
      ))}
      
      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage >= totalPages}
        className="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-700 hover:bg-gray-50 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed transition-colors"
      >
        Next
      </button>
    </div>
  );
}
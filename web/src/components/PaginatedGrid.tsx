import { DataView } from 'primereact/dataview';
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
  layout?: 'grid' | 'list';
}

export function PaginatedGrid<T>({
  data,
  isLoading,
  error,
  onPageChange,
  renderItem,
  emptyState,
  layout = 'grid'
}: PaginatedGridProps<T>) {
  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-8">
        <Spinner size={32} />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-red-600 py-4 text-center">
        Failed to load data. Please try again.
      </div>
    );
  }

  if (!data?.items?.length) {
    return (
      <div className="text-center py-8">
        {emptyState || <span className="text-gray-500">No items found.</span>}
      </div>
    );
  }

  const rows = data.perPage || 10;
  const totalRecords = (data.total ?? 0) || ((data.totalPages ?? 0) * rows);
  const first = ((data.page ?? 1) - 1) * rows;

  const itemTemplate = (item: T, options: any) => {
    return renderItem(item, options.index);
  };

  return (
    <DataView
      value={data.items}
      itemTemplate={itemTemplate}
      layout={layout}
      paginator={totalRecords > rows}
      rows={rows}
      totalRecords={totalRecords}
      lazy
      first={first}
      onPage={(e) => onPageChange?.(e.page + 1)}
      className="app-dataview"
      pt={{
          grid: { className: 'grid gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3' }
      }}
    />
  );
}

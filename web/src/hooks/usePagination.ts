import { useState, useCallback } from 'react';

export interface PaginationParams {
  page?: number;
  size?: number;
}

export function usePagination(initialPage = 0, initialSize = 10) {
  const [page, setPage] = useState(initialPage);
  const [size, setSize] = useState(initialSize);

  const updatePage = useCallback((newPage: number) => {
    setPage(newPage);
  }, []);

  const updateSize = useCallback((newSize: number) => {
    setSize(newSize);
    setPage(0); // Reset to first page when changing size
  }, []);

  const reset = useCallback(() => {
    setPage(initialPage);
    setSize(initialSize);
  }, [initialPage, initialSize]);

  return {
    page,
    size,
    updatePage,
    updateSize,
    reset,
    params: { page: Math.max(0, page - 1), size } as PaginationParams
  };
}
import { useState, useCallback } from 'react';

export interface PaginationParams {
  page?: number;
  size?: number;
}

export function usePagination(initialPage = 1, initialSize = 10) {
  const [page, setPage] = useState(initialPage);
  const [size, setSize] = useState(initialSize);

  const updatePage = useCallback((newPage: number) => {
    setPage(newPage);
  }, []);

  const updateSize = useCallback((newSize: number) => {
    setSize(newSize);
    setPage(1); // Reset to first page when changing size
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
    params: { page, size } as PaginationParams
  };
}
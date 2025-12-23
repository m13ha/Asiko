import { renderHook, act } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { usePagination } from './usePagination';

describe('usePagination', () => {
  it('should initialize with default values', () => {
    const { result } = renderHook(() => usePagination());
    expect(result.current.page).toBe(0);
    expect(result.current.size).toBe(10);
    expect(result.current.params).toEqual({ page: 0, size: 10 });
  });

  it('should initialize with provided values', () => {
    const { result } = renderHook(() => usePagination(2, 25));
    expect(result.current.page).toBe(2);
    expect(result.current.size).toBe(25);
  });

  it('should update page', () => {
    const { result } = renderHook(() => usePagination());
    act(() => {
      result.current.updatePage(5);
    });
    expect(result.current.page).toBe(5);
  });

  it('should update size and reset page', () => {
    const { result } = renderHook(() => usePagination(3, 10));
    act(() => {
      result.current.updateSize(50);
    });
    expect(result.current.size).toBe(50);
    expect(result.current.page).toBe(0);
  });

  it('should reset to initial values', () => {
    const { result } = renderHook(() => usePagination(1, 20));
    act(() => {
      result.current.updatePage(5);
      result.current.updateSize(50);
    });
    
    act(() => {
      result.current.reset();
    });

    expect(result.current.page).toBe(1);
    expect(result.current.size).toBe(20);
  });
});

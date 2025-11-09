export function HandUnderline({ width = 64, height = 10 }: { width?: number; height?: number }) {
  return (
    <svg width={width} height={height} viewBox="0 0 64 10" aria-hidden>
      <path d="M2 6 C 12 2, 20 9, 32 6 S 54 2, 62 6" fill="none" stroke="currentColor" strokeWidth="2" opacity="0.6" />
    </svg>
  );
}


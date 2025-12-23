export function Spinner({ size = 16, className = '' }: { size?: number; className?: string }) {
  return (
    <i 
      className={`pi pi-spin pi-spinner ${className}`} 
      style={{ fontSize: `${size}px` }} 
    />
  );
}
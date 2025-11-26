export function Spinner({ size = 16, className = '' }: { size?: number; className?: string }) {
  return (
    <div 
      className={`animate-spin rounded-full border-2 border-gray-300 border-t-blue-600 ${className}`}
      style={{ width: size, height: size }}
    />
  );
}
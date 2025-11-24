import { ProgressSpinner } from 'primereact/progressspinner';

export function Spinner({ size = 16 }: { size?: number }) {
  return (
    <ProgressSpinner
      style={{ width: size, height: size }}
      strokeWidth="4"
      animationDuration=".8s"
      className="app-spinner"
    />
  );
}

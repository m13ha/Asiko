import type { SkeletonProps } from 'primereact/skeleton';
import { Skeleton as PrimeSkeleton } from 'primereact/skeleton';

export function Skeleton(props: SkeletonProps) {
  return <PrimeSkeleton {...props} className={['app-skeleton', props.className].filter(Boolean).join(' ') || undefined} />;
}

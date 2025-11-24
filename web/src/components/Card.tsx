import styled from 'styled-components';
import type { CardProps as PrimeCardProps } from 'primereact/card';
import { Card as PrimeCard } from 'primereact/card';

export type CardProps = PrimeCardProps;

export function Card({ className, children, ...props }: CardProps) {
  const computedClass = ['app-card', className].filter(Boolean).join(' ');
  return (
    <PrimeCard {...props} className={computedClass || undefined}>
      {children}
    </PrimeCard>
  );
}

export const CardHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
`;

export const CardTitle = styled.h2`
  margin: 0;
  font-size: 18px;
`;

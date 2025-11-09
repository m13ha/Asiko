import styled, { keyframes } from 'styled-components';

const shimmer = keyframes`
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
`;

export const Skeleton = styled.div`
  border-radius: var(--radius);
  background: linear-gradient(
    90deg,
    var(--border) 25%,
    color-mix(in oklab, var(--border) 60%, #fff) 37%,
    var(--border) 63%
  );
  background-size: 400% 100%;
  animation: ${shimmer} 1.2s var(--ease-out) infinite;
`;

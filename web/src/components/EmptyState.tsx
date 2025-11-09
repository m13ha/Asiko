import styled from 'styled-components';

export const EmptyState = styled.div`
  display: grid;
  justify-items: center;
  text-align: center;
  gap: 8px;
  padding: 24px;
  border: 1px dashed color-mix(in oklab, var(--text) 14%, transparent);
  border-radius: var(--radius-lg);
  background: color-mix(in oklab, var(--vintage-cream) 35%, white);
`;

export const EmptyTitle = styled.h3`
  margin: 8px 0 0 0;
  font-size: 18px;
  color: var(--text);
`;

export const EmptyDescription = styled.p`
  margin: 0;
  color: var(--text-muted);
  font-size: 14px;
`;

export const EmptyAction = styled.div`
  margin-top: 8px;
`;


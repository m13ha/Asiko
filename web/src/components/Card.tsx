import styled from 'styled-components';

export const Card = styled.div`
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 16px;
  box-shadow: var(--elev-1);
  transition: transform var(--spring-fast), box-shadow var(--spring-fast);
  &:hover { transform: translateY(-2px); box-shadow: var(--elev-2); }
  &:focus-within { box-shadow: var(--elev-2), 0 0 0 2px color-mix(in oklab, var(--primary) 20%, transparent); }
`;

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

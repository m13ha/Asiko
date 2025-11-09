import styled from 'styled-components';

export const ListItem = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 8px 12px;
  background: var(--bg-elevated);
  transition: background var(--spring-fast), transform var(--spring-fast), box-shadow var(--spring-fast), border-color var(--spring-fast);
  &:hover {
    background: color-mix(in oklab, var(--primary) 3%, var(--bg-elevated));
    transform: translateY(-1px);
    box-shadow: var(--elev-1);
    border-color: color-mix(in oklab, var(--primary) 20%, var(--border));
  }
  &:active { transform: scale(0.99); }
  &:focus-within {
    box-shadow: var(--elev-1), 0 0 0 2px color-mix(in oklab, var(--primary) 20%, transparent);
  }
`;


import styled from 'styled-components';

export const Input = styled.input`
  width: 100%;
  padding: 10px 12px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  background: var(--bg-elevated);
  color: var(--text);
  outline: none;
  transition: border-color var(--spring-fast), box-shadow var(--spring-fast), background var(--spring-fast);
  &:focus {
    border-color: color-mix(in oklab, var(--primary) 35%, var(--border));
    box-shadow: 0 0 0 3px color-mix(in oklab, var(--primary) 20%, transparent);
  }
`;

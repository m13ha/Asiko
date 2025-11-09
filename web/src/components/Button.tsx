import styled from 'styled-components';

export const Button = styled.button<{ variant?: 'primary' | 'ghost' }>`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  background: ${({ variant }) => (variant === 'primary' ? 'var(--primary)' : 'transparent')};
  color: ${({ variant }) => (variant === 'primary' ? 'var(--primary-contrast)' : 'var(--text)')};
  cursor: pointer;
  transition: background var(--spring-fast), opacity var(--spring-fast), transform var(--spring-fast), box-shadow var(--spring-fast);
  box-shadow: none;
  &:hover {
    opacity: 0.95;
    transform: translateY(-1px);
    box-shadow: var(--elev-2);
    background: ${({ variant }) => (variant === 'primary' ? 'color-mix(in oklab, var(--primary) 92%, white)' : 'color-mix(in oklab, var(--primary) 6%, transparent)')};
  }
  &:active { transform: scale(0.98); }
  &:disabled { opacity: 0.6; cursor: not-allowed; }
  &:focus-visible {
    outline: none;
    box-shadow: 0 0 0 2px color-mix(in oklab, var(--primary) 40%, transparent), var(--elev-2);
  }
`;

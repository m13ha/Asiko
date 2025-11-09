import styled from 'styled-components';

export const Field = styled.div`
  display: grid;
  gap: 6px;
  @media (max-width: 768px) { gap: 4px; }
`;

export const FieldLabel = styled.label`
  font-size: 12px;
  color: var(--text);
`;

export const FieldRow = styled.div`
  position: relative;
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  gap: 8px;
  @media (max-width: 640px) {
    grid-template-columns: 1fr;
    gap: 6px;
  }
`;

export const IconSlot = styled.span`
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  color: var(--text-muted);
  pointer-events: none;
`;

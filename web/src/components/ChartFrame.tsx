import styled from 'styled-components';

export const ChartCard = styled.div`
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 12px;
  box-shadow: var(--elev-1);
`;

export const ChartTitle = styled.div`
  font-weight: 600;
  margin-bottom: 8px;
`;

export const ChartArea = styled.div`
  width: 100%;
  height: 220px;
`;

export const ChartLegend = styled.div`
  display: flex;
  gap: 12px;
  align-items: center;
  margin-top: 8px;
  flex-wrap: wrap;
`;

export const LegendItem = styled.div`
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--text-muted);
  font-size: 12px;
`;

export const Swatch = styled.span<{ color: string }>`
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 2px;
  background: ${({ color }) => color};
  box-shadow: 0 0 0 1px rgba(0,0,0,0.06) inset;
`;
